<?php

require_once 'config.php';
require_once 'init.php';
require_once 'Client.php';

$client = new Client($path);
$models = $client->listModels();

$listAll = [];

if(isset($_GET['model'])) {
    $listAll = $client->all($_GET['model']);


    if($listAll === false) {
        header('Location: /login.php');
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

        <ul class="nav">
        <?php foreach($models as $model): ?>
         <li class="nav-item">
             <a class="nav-link" href="/?model=<?php echo $model['name']; ?>"><?php echo $model['name']; ?></a>
         </li>
       <?php endforeach; ?>
        </ul>

        <div class="page-header">
            <h1>Language: {{ .Config.Language }}, Version: {{ .Config.Version }} </h1>
        </div>

        <?php /* <p>Query API to: <?php echo $path; ?></p> */ ?>





         <?php if(!empty($listAll)) :?>


         <div>
         <a href="/?model=<?php echo $_GET['model']; ?>&a=new">Add new model</a>
         </div>

            <table class="table">
         <?php foreach($listAll as $n => $item): ?>


            <?php if($n == 0): ?>
            <thead>
            <tr>

                <?php foreach(array_keys(get_object_vars ($item)) as $key): ?>
                <th>
                                <?php echo $key; ?>
                </th>
                <?php endforeach ?>

                <th>
                    Actions
                </th>
            </tr>
            </thead>
            <tbody>
            <?php endif;  ?>



            <tr>
                <?php foreach($item as $field): ?>
                <td>
                    <?php echo $field; ?>
                </td>
                <?php endforeach ?>
                <td>
                    <a href="/?model=<?php echo $_GET['model']; ?>&a=edit&id=<?php echo $item->id; ?>">Edit</a>
                    <a href="/?model=<?php echo $_GET['model']; ?>&a=remove&id=<?php echo $item->id; ?>" onclick="return confirm('?');">Delete</a>
                </td>
             </tr>



         <?php endforeach ?>

                </tbody>
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