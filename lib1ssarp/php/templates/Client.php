<?php

class Client {



    private $baseUrl;

    private $token;

    public function __construct($baseUrl){
        $this->baseUrl = $baseUrl;
    }


    public function listModels() {
        $models = [
            {{ range .Config.Models }}
                [
                    'name'  => '{{.Name}}',
                ],
            {{ end }}
        ];

        return  $models;
    }


    public function all($model) {
        return $this->query('api/'.$model);
    }


    private function query($uri, $method = 'GET', $body = null) {
        $opts = [
            "http" => [
                "method" => $method,
                "header" => "Accept-language: en\r\n" .
                    "Authorization: Token {$this->getToken()}\r\n" .
                    "Content-type: application/x-www-form-urlencoded\r\n",
                'content' => json_encode($body),
            ]
        ];

        $url = $this->baseUrl . '/' . $uri;

        $context = stream_context_create($opts);
        $response = file_get_contents($url, false, $context);

        $data = json_decode($response);

        if(!$data) {
            return false;
        }

      //  var_dump($response);
        return $data;
    }


    public function openTokenIsEmpty( $force = false) {
        if(empty($this->getToken()) || $force) {

            $sold = md5(time() . __METHOD__);
            $data = $this->query('session', 'POST', [
                'sold'  => $sold,
                'token' => md5($this->getServerToken() . $sold . 'POST')
            ]);


            if($data && $data->session) {
                $this->token = $data->session;
                $_SESSION['token_session_'] = $this->token;
            }
        }

    }

    public function getServerToken() {
        return 'xsxksmkxmskxmskxmskxmskx';
    }


    public function getToken() {
        return isset($_SESSION['token_session_']) ? $_SESSION['token_session_'] : $this->token;
    }
}