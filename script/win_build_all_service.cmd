SET ROOT=%cd%
mkdir %ROOT%\..\bin\
cd ..\cmd\open_im_api\ && go build  && move open_im_api.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_cms_api\ && go build  && move open_im_cms_api.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_demo\ && go build  && move open_im_demo.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_msg_gateway\ && go build  && move open_im_msg_gateway.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_msg_transfer\ && go build  && move open_im_msg_transfer.exe %ROOT%\..\bin\
cd ..\..\cmd\open_im_push\ && go build  && move open_im_push.exe %ROOT%\..\bin\
cd ..\..\cmd\rpc\open_im_admin_cms\&& go build  && move open_im_admin_cms.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_auth\&& go build  && move open_im_auth.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_cache\&& go build  && move open_im_cache.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_conversation\&& go build  && move open_im_conversation.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_friend\&& go build  && move open_im_friend.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_group\&& go build  && move open_im_group.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_message_cms\&& go build  && move open_im_message_cms.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_msg\&& go build  && move open_im_msg.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_office\&& go build  && move open_im_office.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_organization\&& go build  && move open_im_organization.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_statistics\&& go build  && move open_im_statistics.exe %ROOT%\..\bin\
cd ..\..\..\cmd\rpc\open_im_user\&& go build  && move open_im_user.exe %ROOT%\..\bin\

cd %ROOT%
cd ..\cmd\rpc\open_im_task\&& go build  && move open_im_task.exe %ROOT%\..\bin\

cd %ROOT%
cd ..\cmd\rpc\open_im_order\&& go build  && move open_im_order.exe %ROOT%\..\bin\

cd %ROOT%
cd ..\cmd\rpc\open_im_ens\&& go build  && move open_im_ens.exe %ROOT%\..\bin\

cd %ROOT%
cd ..\cmd\rpc\open_im_web3\&& go build  && move open_im_web3.exe %ROOT%\..\bin\

cd ..\..\..\cmd\Open-IM-SDK-Core\ws_wrapper\cmd\&& go build  open_im_sdk_server.go && move open_im_sdk_server.exe %ROOT%\..\bin\
cd %ROOT%