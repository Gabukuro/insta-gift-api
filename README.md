# insta-gift-api

## Visão Geral

A `insta-gift-api` é uma API escrita em Golang v1.20 para fornecer funcionalidades específicas. O projeto utiliza o Makefile para facilitar o gerenciamento de tarefas comuns, e inclui comandos para inicializar, instalar dependências, executar, testar e realizar outras operações úteis.

## Pré-requisitos

Certifique-se de ter o Golang v1.20 ou superior instalado em sua máquina.

## Comandos do Makefile

### Inicializar o Projeto e Instalar Dependências

```bash
make init
make install
```
Este comando inicializa o projeto e instala as dependências necessárias.

### Rodar a Aplicação em Desenvolvimento
```bash
make run-dev
```
Este comando inicia a aplicação em modo de desenvolvimento usando o `nodemon` para hot reload.

Para mais detalhes sobre os comandos disponíveis, consulte o próprio Makefile no projeto.
