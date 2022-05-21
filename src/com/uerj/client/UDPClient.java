package com.uerj.client;

import java.io.*;
import java.net.*;
import java.util.Arrays;
import java.util.Timer;
import java.util.TimerTask;

class UDPClient {
    public static final String DIVISORIA = "=================================================\n";
    public static final String MSG_START = "Sistema Cliente UDP Iniciado\n";
    public static final String MSG_BOAS_VINDAS = DIVISORIA + MSG_START + DIVISORIA;

    public static final String LOCALHOST = "LocalHost";

    public static long stopWatchStartTime = 0;
    public static long stopWatchStopTime = 0;
    public static boolean stopWatchRunning = false;

    public static void main(String args[]) throws IOException {

        boolean validIP = false;

        String ip = null;
        String tipo = null;
        String valor = null;
        InetAddress IPAddress = null;

        BufferedReader buffer = new BufferedReader(new InputStreamReader(System.in));

        System.out.println(MSG_BOAS_VINDAS);

        do {
            try {

                System.out.println("Digite o IP ou deixe em branco para usar o padrao(LocalHost): ");

                ip = buffer.readLine();

                if (ip.equals("")){ ip = LOCALHOST; }

                IPAddress = InetAddress.getByName(ip);

                validIP = true;
            }catch (IOException e) {
                System.out.println("Endereco Invalido");
            }
        } while (!validIP);

        System.out.println("Digite a porta desejada: ");

        int porta = Integer.parseInt(buffer.readLine());

        do {
            System.out.println("Digite o tipo da mensagem: ");

            tipo = buffer.readLine();

        } while (!validaTipo(tipo));

        do {
            System.out.println("Digite o valor da mensagem: ");

            valor = buffer.readLine();
        }while (!validaValor(valor,Tipo.valueOf(tipo.toUpperCase())));


        DatagramSocket clientSocket = new DatagramSocket();

        byte[] sendData = new byte[1024];
        byte[] receiveData = new byte[1024];

        String sentence = criaJson(tipo,  valor);

        System.out.println("Retorno: \n" + sentence + "\nIP: " + ip + "\nPorta: " + porta +"\n");

        sendData = sentence.getBytes();

        DatagramPacket sendPacket = new DatagramPacket(sendData, sendData.length, IPAddress, porta);

        clientSocket.send(sendPacket);
        //iniciar timer
        startTimer();
        System.out.println("Mensagem Enviada com sucesso!");

        timeOut();

        DatagramPacket receivePacket = new DatagramPacket(receiveData, receiveData.length);

        clientSocket.receive(receivePacket);
        //parar timer e exibir
        stopTimer();

        String modifiedSentence = new String(receivePacket.getData());

        System.out.println("FROM SERVER: " + modifiedSentence);
        System.out.println("Tempo decorrido: " + getElapsedMilliseconds() + "MiliSegundos");

        clientSocket.close();

    }

    public static String criaJson(String tipo, String valor){
        MensagemCliente mensagemCliente = new MensagemCliente(tipo, valor);

        return mensagemCliente.toString();
    }

    public static boolean validaTipo(String tipo){

        try {
            Tipo.valueOf(tipo.toUpperCase());
        }catch (IllegalArgumentException e){
            System.out.println("Tipo Invalido!\n A mensagem deve ser de um dos seguintes tipos: " + Arrays.toString(Tipo.values()));
            return false;
        }

        return true;

    }

    public static boolean validaValor(String valor, Tipo tipo){

        String errorMsg = "Valor invalido!\nPor favor digite um ";

        switch (tipo){
            case INT:
                try {
                    Integer.parseInt(valor);
                }catch (NumberFormatException e){
                    System.out.println(errorMsg + "inteiro");
                    return false;
                }
                return true;
            case CHAR:
                if (valor.length() > 1){
                    System.out.println(errorMsg + "caractere");
                    return false;
                }

        }

        return true;

    }

    public static void startTimer() {
        stopWatchStartTime = System.nanoTime();
        stopWatchRunning = true;
    }


    public static void stopTimer() {
        stopWatchStopTime = System.nanoTime();
        stopWatchRunning = false;
    }


    public static long getElapsedMilliseconds() {
        long elapsedTime;

        if (stopWatchRunning)
            elapsedTime = (System.nanoTime() - stopWatchStartTime);
        else
            elapsedTime = (stopWatchStopTime - stopWatchStartTime);

        long nanoSecondsPerMillisecond = 1000000;
        return elapsedTime / nanoSecondsPerMillisecond;
    }

    public static void timeOut(){

        long timeout = 60000;
        TimerTask task = new TimerTask() {
            @Override
            public void run() {
                System.out.println("\nTime Out, Encerrando a aplicacao...");
                System.exit(1);
            }
        };

        Timer timer = new Timer();
        timer.schedule(task, timeout);
    }


}

