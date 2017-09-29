<?php

$host = '{{ .Config.Server.Host }}';
$port = '{{ .Config.Server.Port }}';
if(empty($host)) {
    $host = 'localhost';

}
$path = 'http://' .$host . ':'.$port;
