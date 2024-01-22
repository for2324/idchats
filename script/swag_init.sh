# 生成swag文档 需在项目根目录执行
if [ "$(uname)" == "Darwin" ]; then
 PROJECTROOT=/Users/apple/Documents/Open-IM-Server
 WorkProject=/Users/apple/Documents/gosrc/src/Open-IM-Server
 /bin/cp -rf ../pkg/base_info $PROJECTROOT/pkg/
 /bin/cp -rf ../pkg/common/db $PROJECTROOT/pkg/common/
 /bin/cp -rf ../pkg/proto $PROJECTROOT/pkg/
 /bin/cp -r ../internal/api   $PROJECTROOT/internal/
 /bin/cp -r ../cmd/open_im_api/main.go   $PROJECTROOT/cmd/open_im_api/main.go
 cd $PROJECTROOT
elif [[ "$(uname)" == "MINGW"* ]]; then
  SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
  cd $SCRIPT_DIR/../ # 进入工作目录
fi

swag init --parseVendor --parseInternal --parseDependency --parseDepth 100   -o ./cmd/open_im_api/docs  -g ./cmd/open_im_api/main.go
if [ "$(uname)" == "Darwin" ]; then
 /bin/cp -rf  $PROJECTROOT/cmd/open_im_api/docs   $WorkProject/cmd/open_im_api/
 cd $WorkProject/script
fi
