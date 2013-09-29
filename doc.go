/*
You can use this tool to log http request, including request headers and values in POST/GET form.

It has two kinds of methods to log the request, local and remote.
Local logger need a function like log.Log whose type is func(...interface{}),
while remote logger will send the log by socket to another server.

The format of logger is designed accourding to the internal tool in NetEase Youdao.

*/
package requestlog
