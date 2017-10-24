<?php

require_once 'config.php';
require_once 'init.php';
require_once 'Client.php';

$client = new Client($path);



if(isset($_POST['login']) && isset($_POST['password']) && $_POST['login'] == 'admin' && $_POST['password'] == 'admin') {
    $client->openTokenIsEmpty(true);

    if($client->getToken()) {
         header('Location: /');
         exit(0);
    }
}

?>

<!DOCTYPE html>
    <html lang="en">
    <head>

        <meta charset="utf-8">
        <title>PHP client</title>

        <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/css/bootstrap.min.css" integrity="sha384-/Y6pD6FV/Vv2HJnA6t+vslU6fwYXjCFtcEpHbNJ0lyAFsXTsjBbfaDjzALeQsN6M" crossorigin="anonymous">
        <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
        </head>
    <body>


     <div class="container">

     <!-- form login -->
 <form role="form" method="post">
                        <fieldset>
                            <div class="form-group">
                                <input class="form-control" placeholder="login" name="login" type="login" autofocus>
                            </div>
                            <div class="form-group">
                                <input class="form-control" placeholder="Password" name="password" type="password" value="">
                            </div>

                            <button type="submit" class="btn btn-success btn-block">Login</button>
                            <p>admin - admin</p>
                        </fieldset>
                    </form>

     <!-- /form login -->

         </div>

         <footer class="footer">
            <div class="container">
            <p class="text-muted">Place sticky footer content here.</p>
          </div>
        </footer>


            <script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
            <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.11.0/umd/popper.min.js" integrity="sha384-b/U6ypiBEHpOf/4+1nzFpr53nxSS+GLCkfwBdFNTxtclqqenISfwAzpKaMNFNmj4" crossorigin="anonymous"></script>
            <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/js/bootstrap.min.js" integrity="sha384-h0AbiXch4ZDo7tp9hKZ4TsHbi047NrKGLO3SEJAg45jXxnGIfYzk4Si90RDIqNm1" crossorigin="anonymous"></script>

        </body>
    </html>