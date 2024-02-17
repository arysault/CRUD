package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// CriarUsuario cria um usuario no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da Requisição"))
		return
	}

	var usuario usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Nãofoi possivel converter usuario para Struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados")) //Erro interno do Servidor
		return
	}
	defer db.Close()

	stmt, erro := db.Prepare("insert into usuarios(nome, email) values (? ,?)")
	if erro != nil {
		w.Write([]byte("Falha na preparação da query: " + erro.Error())) //Erro na query SQL
		return
	}
	defer stmt.Close()

	resultado, erro := stmt.Exec(usuario.Nome, usuario.Email)
	if erro != nil {
		w.Write([]byte("Falha no exec")) //Erro no exec
		return
	}

	idInserido, erro := resultado.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao obter Id inserido")) //Erro ao obter novo Id
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso, Id: %d", idInserido)))

}

// BuscarUsuarios Busca todos os usuarios do banco de dados
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o Banco de Dados!"))
		return
	}
	defer db.Close()

	linhas, erro := db.Query("select * from usuarios")
	if erro != nil {
		w.Write([]byte("Erro ao buscar ususarios"))
		return
	}
	defer linhas.Close()

	var usuarios []usuario

	for linhas.Next() {
		var usuario usuario
		if erro := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro na leitura dos dados"))
			return
		}

		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil {
		w.Write([]byte("Erro ao converter para json"))
		return
	}

}

// BuscarUsuario traz apenas um usuário no banco de dados
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Não foi possivel converter o paramentro para inteiro!"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados: " + erro.Error()))
		return
	}

	linha, erro := db.Query("select * from usuarios where id = ?", ID)
	if erro != nil {
		w.Write([]byte("Erro na busca do usuario: " + erro.Error()))
		return
	}

	var usuario usuario
	if linha.Next() {
		if erro := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear Usuario"))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		w.Write([]byte("Erro ao transformar o usuario em json!"))
	}

}

// AtualizarUsuario altera usuário no banco de dados
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {

	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao ler parametro id"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Erro ao ler o corpo da requisição"))
		return
	}

	var usuario usuario

	if erro := json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao decodificar os dados"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro na conexão com o banco de dados"))
		return
	}
	defer db.Close()

	stmt, erro := db.Prepare("update usuarios set nome = ?, email = ?  where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar statement!"))
		return
	}
	defer stmt.Close()

	if _, erro := stmt.Exec(usuario.Nome, usuario.Email, ID); erro != nil {
		w.Write([]byte("Erro ao atualizar o usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// DeletarUsuario remove usuario do bando de dados
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametro := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametro["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter paramentro"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar ao banco de dados"))
		return
	}
	defer db.Close()

	stmt, erro := db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar Statement!"))
		return
	}
	defer stmt.Close()

	if _, erro := stmt.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao deletar Usuário!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
