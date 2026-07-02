# Seguranca em operacoes PIX

Operacoes PIX sao consideradas acoes criticas porque movimentam dinheiro.

Antes da execucao, o agente deve solicitar confirmacao explicita do usuario.

A LLM pode propor a operacao, mas a confirmacao deve ser controlada pelo backend e nunca inferida pela LLM.
