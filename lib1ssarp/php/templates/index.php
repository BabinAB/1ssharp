<?php

//TODO move to config file
$path = 'http://localhost:{{ .Config.Server.Port }}';

?>
<html>
    <head></head>
    <body>
        <h1>Language: {{ .Config.Language }}, Version: {{ .Config.Version }} </h1>

        <p>Query API to: <?php echo $path; ?></p>
    </body>
</html>