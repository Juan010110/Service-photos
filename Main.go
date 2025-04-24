package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Abrir el archivo de texto con las URLs
	file, err := os.Open("hola.txt")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Crear directorio para guardar las imágenes si no existe
	err = os.MkdirAll("imagenes_descargadas", 0755)
	if err != nil {
		fmt.Println("Error al crear el directorio de imágenes:", err)
		return
	}

	// Leer el archivo línea por línea
	scanner := bufio.NewScanner(file)
	lineNum := 0
	descargadas := 0

	for scanner.Scan() {
		lineNum++
		url := strings.TrimSpace(scanner.Text())
		
		// Omitir líneas vacías
		if url == "" {
			continue
		}

		// Verificar si la URL es de Imgur
		if !strings.Contains(strings.ToLower(url), "imgur") {
			fmt.Printf("Línea %d: %s - No es una URL de Imgur, omitiendo...\n", lineNum, url)
			continue
		}

		descargadas++
		fmt.Printf("Descargando imagen %d de Imgur: %s\n", descargadas, url)
		
		// Obtener el nombre del archivo de la URL
		fileName := filepath.Base(url)
		if fileName == "" || fileName == "." {
			fileName = fmt.Sprintf("imgur_%d.jpg", descargadas)
		}
		
		// Descargar la imagen
		err := descargarImagen(url, filepath.Join("imagenes_descargadas", fileName))
		if err != nil {
			fmt.Printf("Error al descargar %s: %s\n", url, err)
		} else {
			fmt.Printf("Imagen de Imgur guardada como: %s\n", fileName)
		}
		
		// Esperar 10 segundos antes de la siguiente descarga
		fmt.Println("Esperando 10 segundos...")
		time.Sleep(10 * time.Second)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}

	fmt.Printf("Proceso completado. Se descargaron %d imágenes de Imgur.\n", descargadas)
}

func descargarImagen(url, rutaDestino string) error {
	// Crear una solicitud HTTP GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	
	// Añadir cabeceras para simular un navegador
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://imgur.com/")
	
	// Realizar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificar que la respuesta sea exitosa
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("código de estado HTTP no válido: %d", resp.StatusCode)
	}

	// Crear el archivo de destino
	out, err := os.Create(rutaDestino)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copiar los datos de la respuesta al archivo
	_, err = io.Copy(out, resp.Body)
	return err
}