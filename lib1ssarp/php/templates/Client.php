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
        $this->openTokenIsEmpty();
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
       // var_dump($response);
        return $data;
    }


    public function openTokenIsEmpty() {
        if(is_null($this->token)) {

            $sold = md5(time() . __METHOD__);
            $data = $this->query('session', 'POST', [
                'sold'  => $sold,
                'token' => md5('xsxksmkxmskxmskxmskxmskx' . $sold . 'POST')
            ]);


            if($data && $data->session) {
                $this->token = $data->session;
            }
        }

    }

    public function getToken() {
        return $this->token;
    }
}