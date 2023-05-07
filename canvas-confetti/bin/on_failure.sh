# 작업이 실패했을때 실행되는 스크립트 스크린샷을 찍어 파일로 저장하며 이를 transfer.sh를 통해 업로드 실행

green=`tput setaf 2` #setaf -> ANSI Escape Codes 중 색상 지정 코드 2는 녹색을 뜻함
reset=`tput sgr0` # sgr0 -> 출력할 문자열의 스타일 리셋을 위한 코드 reset은 문자열의 색상 리셋을 위한 코드가 담김
#tput은 ANSI Escape Codes를 출력하거나 터미널의 속성 값을 리턴하는 역할

line_break () {
  echo --------------------------------------
}

print_green () {
  echo "${green}$@${reset}"
}

upload_file () {
  filename=$1
  urlname=${filename// /_}

  downloadurl=`curl -sS --upload-file "./$filename" https://transfer.sh/$urlname`

  echo image \"$filename\"
  print_green "  uploaded to: $downloadurl"
}

find_files () {
  cd shots

  echo list of files present:
  ls -l
  line_break

  for i in *.png;do upload_file "$i";done
}

find_files
