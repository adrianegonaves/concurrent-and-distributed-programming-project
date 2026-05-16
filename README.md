# 🎭 Gotion: Hybrid Emotion Recognition System

**Trabalho Final | Programação Concorrente e Distribuída**  
Este projeto consiste num sistema híbrido de reconhecimento de expressões faciais (FER), utilizando **Go** para processamento paralelo massivo de dados e **Python/TensorFlow** para a rede neural convolucional.

---

## 👥 Equipa
*   **Adriane Gonçalves** - 240000004
*   **Bruno Hortelão** - 240001083

---
---

## ✅ Checklist de Implementação

### 1. 🏎️ Go Pre-processing Pipeline (Paralelismo Real)
- [ ] **Worker Pool Pattern:** Implementação de múltiplos operários para processamento simultâneo utilizando todos os núcleos da CPU.
- [ ] **Image Resizer:** Redimensionamento uniforme para 48x48 pixels utilizando algoritmos de interpolação (ex: Lanczos).
- [ ] **Grayscale Converter:** Garantia de canal único (1-channel) para consistência da entrada da rede neural.
- [ ] **Concurrent File Walker:** Leitura eficiente e rápida da estrutura de pastas do dataset FER-2013.

### 2. 🧠 Python Neural Network (Deep Learning)
- [ ] **Data Pipeline:** Consumo otimizado dos dados processados via `tf.data` com suporte a *prefetch* e *caching*.
- [ ] **CNN Architecture:** Construção de camadas Convolucionais, Batch Normalization e Dropout para extração de características.
- [ ] **Model Training:** Implementação de callbacks como *Early Stopping* e *Learning Rate Reduction*.
- [ ] **Evaluation Metrics:** Geração de Matriz de Confusão e Relatórios de Classificação (Precision, Recall, F1).

### 3. 📂 Data Management
- [ ] **Raw Data Handling:** Gestão e organização do dataset original (bruto).
- [ ] **Clean Data Export:** Exportação para uma estrutura de diretórios otimizada em `data/processed`.
- [ ] **Model Serialization:** Guardar o modelo treinado em formato `.h5` ou `.keras` para inferência futura.

---

## 🛠️ Tecnologias Utilizadas

| Tecnologia | Função |
| :--- | :--- |
| **Go (Golang)** | Processamento Paralelo e Engenharia de Dados (CPU-bound tasks) |
| **Python** | Treino de Deep Learning e Pesquisa de IA |
| **TensorFlow/Keras** | Framework principal para a Rede Neural Convolucional |
| **Imaging (Go Library)** | Biblioteca de manipulação rápida de ficheiros de imagem |
| **Matplotlib/Seaborn** | Visualização de dados e métricas de desempenho |

[Golang com struct, funções e metodos](https://tomelin-tech.medium.com/golang-com-struct-funções-e-metodos-698a25b6221a)
[](https://gobyexample.com/channels)
[](https://pkg.go.dev/github.com/disintegration/imaging#section-readme)
---

