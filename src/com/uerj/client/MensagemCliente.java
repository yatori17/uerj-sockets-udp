package com.uerj.client;

public class MensagemCliente {
    Tipo tipo;
    String valor;

    public MensagemCliente(String tipo, String valor) {
        this.tipo = Tipo.valueOf(tipo.toUpperCase());
        this.valor = valor;
    }

    @Override
    public String toString() {

        switch (tipo){

            case CHAR:
                return "{" +
                        "\"tipo\":" + "\"" + tipo.toString().toLowerCase() + "\"" +
                        ", \"valor\":" + "\'" + valor + '\'' +
                        '}';

            case STRING:
                return "{" +
                        "\"tipo\":" + "\"" + tipo.toString().toLowerCase() + "\"" +
                        ", \"valor\":" + "\"" + valor + '\"' +
                        '}';
            case INT:
                return "{" +
                        "\"tipo\":" + "\"" + tipo.toString().toLowerCase() + "\"" +
                        ", \"valor\":" + valor +
                        '}';
        }

        return "Tipo invalido";
    }


}
