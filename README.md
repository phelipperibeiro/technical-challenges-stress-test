# Stress Test

## Como executar

#### 1) Construa a imagem docker
```bash
docker build -t technical-challenges-stress-test .
```

#### 2) Após a imagem ser construída, podemos executar a imagem com docker run, por exemplo:

```bash
docker run technical-challenges-stress-test --url=http://globo.com --requests=1000 --concurrency=10
```

#### 2.1) Os parâmetros são:
- url: A URL a ser testada
- requests: O número de solicitações a serem feitas
- concurrency: O número de solicitações a serem feitas simultaneamente


#### 3) A saída será algo como:
```bash
Report:
Total requests: 1000
Time taken: 1.000s
Status code distribution:
  [200] 1000 requests
```
