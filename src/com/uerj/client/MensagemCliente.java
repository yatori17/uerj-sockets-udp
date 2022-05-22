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

            case STRING:
                return "{" +
                        "\"tipo\":" + "\"" + tipo.toString().toLowerCase() + "\"" +
                        ", \"val\":" + "\"" + valor + '\"' +
                        '}';
            case INT:
                return "{" +
                        "\"tipo\":" + "\"" + tipo.toString().toLowerCase() + "\"" +
                        ", \"val\":" + valor +
                        '}';
        }

        return "Tipo invalido";
    }


}
