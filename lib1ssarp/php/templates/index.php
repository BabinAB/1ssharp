<?php

require_once 'config.php';
require_once 'Client.php';

$client = new Client($path);
$models = $client->listModels();

$listAll = [];

if(isset($_GET['model'])) {
    $listAll = $client->all($_GET['model']);
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


     <nav class="navbar navbar-default">
             <div class="container-fluid">
               <div class="navbar-header">
                 <a class="navbar-brand" href="#">Project name</a>
               </div>

             </div>
           </nav>

        <div class="page-header">
            <h1>Language: {{ .Config.Language }}, Version: {{ .Config.Version }} </h1>
        </div>

        <p>Query API to: <?php echo $path; ?></p>

         <?php foreach($models as $model): ?>
            <a href="/?model=<?php echo $model['name']; ?>"><?php echo $model['name']; ?></a>

         <?php endforeach; ?>



         <?php if(!empty($listAll)) :?>
            <table class="table">
         <?php foreach($listAll as $item): ?>

            <tr>
                <?php foreach($item as $field): ?>
                <td>
                    <?php echo $field; ?>
                </td>
                <?php endforeach ?>
             </tr>
         <?php endforeach ?>
            </table>
         <?php endif ?>


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