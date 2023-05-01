use std::io::Write;

use anyhow::{bail, Result};
use indoc::writedoc;
use turbo_binding::{
    turbo::{
        tasks::{primitives::StringVc, TryJoinIterExt, Value},
        tasks_fs::rope::RopeBuilder,
    },
    turbopack::{
        core::{
            asset::{Asset, AssetContentVc, AssetVc, AssetsVc},
            chunk::{
                availability_info::AvailabilityInfo, ChunkGroupReferenceVc, ChunkItem, ChunkItemVc,
                ChunkVc, ChunkableAsset, ChunkableAssetVc, ChunkingContext, ChunkingContextVc,
            },
            ident::AssetIdentVc,
            reference::AssetReferencesVc,
        },
        dev::ChunkData,
        ecmascript::{
            chunk::{
                EcmascriptChunkItem, EcmascriptChunkItemContent, EcmascriptChunkItemContentVc,
                EcmascriptChunkItemVc, EcmascriptChunkPlaceable, EcmascriptChunkPlaceableVc,
                EcmascriptChunkVc, EcmascriptChunkingContextVc, EcmascriptExports,
                EcmascriptExportsVc,
            },
            utils::StringifyJs,
        },
    },
};

#[turbo_tasks::function]
fn modifier() -> StringVc {
    StringVc::cell("chunks".to_string())
}

#[turbo_tasks::value(shared)]
pub struct WithChunksAsset {
    pub asset: EcmascriptChunkPlaceableVc,
    pub chunking_context: ChunkingContextVc,
}

#[turbo_tasks::value_impl]
impl Asset for WithChunksAsset {
    #[turbo_tasks::function]
    fn ident(&self) -> AssetIdentVc {
        self.asset.ident().with_modifier(modifier())
    }

    #[turbo_tasks::function]
    fn content(&self) -> AssetContentVc {
        unimplemented!()
    }

    #[turbo_tasks::function]
    async fn references(self_vc: WithChunksAssetVc) -> Result<AssetReferencesVc> {
        let this = self_vc.await?;
        let entry_chunk = self_vc.entry_chunk();

        Ok(AssetReferencesVc::cell(vec![ChunkGroupReferenceVc::new(
            this.chunking_context,
            entry_chunk,
        )
        .into()]))
    }
}

#[turbo_tasks::value_impl]
impl ChunkableAsset for WithChunksAsset {
    #[turbo_tasks::function]
    fn as_chunk(
        self_vc: WithChunksAssetVc,
        context: ChunkingContextVc,
        availability_info: Value<AvailabilityInfo>,
    ) -> ChunkVc {
        EcmascriptChunkVc::new(
            context,
            self_vc.as_ecmascript_chunk_placeable(),
            availability_info,
        )
        .into()
    }
}

#[turbo_tasks::value_impl]
impl EcmascriptChunkPlaceable for WithChunksAsset {
    #[turbo_tasks::function]
    async fn as_chunk_item(
        self_vc: WithChunksAssetVc,
        context: EcmascriptChunkingContextVc,
    ) -> Result<EcmascriptChunkItemVc> {
        Ok(WithChunksChunkItem {
            context,
            inner: self_vc,
        }
        .cell()
        .into())
    }

    #[turbo_tasks::function]
    fn get_exports(&self) -> EcmascriptExportsVc {
        // TODO This should be EsmExports
        EcmascriptExports::Value.cell()
    }
}

#[turbo_tasks::value_impl]
impl WithChunksAssetVc {
    #[turbo_tasks::function]
    async fn entry_chunk(self) -> Result<ChunkVc> {
        let this = self.await?;
        Ok(this.asset.as_root_chunk(this.chunking_context))
    }

    #[turbo_tasks::function]
    async fn chunks(self) -> Result<AssetsVc> {
        let this = self.await?;
        Ok(this.chunking_context.chunk_group(self.entry_chunk()))
    }
}

#[turbo_tasks::value]
struct WithChunksChunkItem {
    context: EcmascriptChunkingContextVc,
    inner: WithChunksAssetVc,
}

#[turbo_tasks::value_impl]
impl EcmascriptChunkItem for WithChunksChunkItem {
    #[turbo_tasks::function]
    fn chunking_context(&self) -> EcmascriptChunkingContextVc {
        self.context
    }

    #[turbo_tasks::function]
    async fn content(&self) -> Result<EcmascriptChunkItemContentVc> {
        let inner = self.inner.await?;
        let Some(inner_chunking_context) = EcmascriptChunkingContextVc::resolve_from(inner.chunking_context).await? else {
            bail!("the chunking context is not an EcmascriptChunkingContextVc");
        };

        let chunks = self.inner.chunks().await?;
        let server_root = inner_chunking_context.output_root().await?;

        let chunks_data = chunks
            .iter()
            .map(|&chunk| ChunkData::from_asset(&server_root, chunk))
            .try_join()
            .await?
            .into_iter()
            .flatten()
            .collect::<Vec<_>>();

        let module_id = &*inner
            .asset
            .as_chunk_item(inner_chunking_context)
            .id()
            .await?;

        let mut code = RopeBuilder::default();

        writedoc!(
            code,
            r#"
            __turbopack_esm__({{
                default: () => {},
                chunks: () => chunks,
            }});
            const chunks = {:#};
            "#,
            StringifyJs(&module_id),
            StringifyJs(&chunks_data),
        )?;

        Ok(EcmascriptChunkItemContent {
            inner_code: code.build(),
            ..Default::default()
        }
        .cell())
    }
}

#[turbo_tasks::value_impl]
impl ChunkItem for WithChunksChunkItem {
    #[turbo_tasks::function]
    fn asset_ident(&self) -> AssetIdentVc {
        self.inner.ident()
    }

    #[turbo_tasks::function]
    fn references(&self) -> AssetReferencesVc {
        self.inner.references()
    }
}
