import os
import shutil
import kagglehub

print("Iniciando o download do Animal Image Dataset via KaggleHub")

# Baixa o dataset de acordo com o Kagglehub
path = kagglehub.dataset_download('iamsouravbanerjee/animal-image-dataset-90-different-animals')
print("Download concluido", path)


diretorio_do_script = os.path.dirname(os.path.abspath(__file__))

# Força o destino a ser a pasta 'origin' ao lado do arquivo 'teste.py'
destino_origin = os.path.join(diretorio_do_script, "data/origin")
os.makedirs(destino_origin, exist_ok=True)
# =============================================================

# Mapeia a estrutura interna do dataset baixado
pasta_das_imagens = os.path.join(path, "animals", "animals")

if not os.path.exists(pasta_das_imagens):
    pasta_das_imagens = os.path.join(path, "animals")
if not os.path.exists(pasta_das_imagens):
    pasta_das_imagens = path

print(f"Copiando pastas de animais para o repositório em: {destino_origin}")

# Varre e copia as pastas reais das espécies para dentro de 'origin'
if os.path.exists(pasta_das_imagens):
    for item in os.listdir(pasta_das_imagens):
        item_path = os.path.join(pasta_das_imagens, item)
        
        if os.path.isdir(item_path):
            dest_file = os.path.join(destino_origin, item)
            if not os.path.exists(dest_file):
                shutil.copytree(item_path, dest_file)
                print(f"-> Pasta [{item}] copiada com sucesso.")

print("\n Concluído!")