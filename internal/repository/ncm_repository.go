package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/IsaqueAmorim/codexfiscal/internal/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type NCMRepository interface {
	CreateNCM(ncm *model.NCM) error
	UpdateNCM(ncm *model.NCM) error
	DeleteNCM(id string) error
	GetAllNCMs() ([]*model.NCM, error)
	GetNCMByCode(code string) (*model.NCM, error)
	GetNCMByID(id string) (*model.NCM, error)
	GetNCMByText(text string) (*model.NCM, error)
	GetNCMSByCodes(codes []string) ([]*model.NCM, error)
	GetNCMSByText(text string) ([]*model.NCM, error)
	BulkInsertNCMs(ncms []*model.NCM) error
}

type ncmRepository struct {
	db *sql.DB
}

func NewNCMRepository(db *sql.DB) NCMRepository {
	return &ncmRepository{
		db: db,
	}
}

func (r *ncmRepository) CreateNCM(ncm *model.NCM) error {
	query := `
		INSERT INTO ncm (id, code, code_no_symbols, description, initial_date, final_date, type_year_ini, number_ato_ini, year_ato_ini)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query, uuid.New().String(), ncm.Code, ncm.CodeNoSimbols, ncm.Description,
		ncm.InitialDate, ncm.FinalDate, ncm.TypeYearIni, ncm.NumberAtoIni, ncm.YearAtoIni)

	if err != nil {
		return fmt.Errorf("erro ao criar NCM: %w", err)
	}

	return nil
}

func (r *ncmRepository) UpdateNCM(ncm *model.NCM) error {
	query := `
		UPDATE ncm 
		SET code = $2, code_no_symbols = $3, description = $4, initial_date = $5, 
			final_date = $6, type_year_ini = $7, number_ato_ini = $8, year_ato_ini = $9
		WHERE id = $1
	`

	result, err := r.db.Exec(query, ncm.ID, ncm.Code, ncm.CodeNoSimbols, ncm.Description,
		ncm.InitialDate, ncm.FinalDate, ncm.TypeYearIni, ncm.NumberAtoIni, ncm.YearAtoIni)

	if err != nil {
		return fmt.Errorf("erro ao atualizar NCM: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("NCM com ID %s não encontrado", ncm.ID)
	}

	return nil
}

func (r *ncmRepository) DeleteNCM(id string) error {
	query := `DELETE FROM ncm WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar NCM: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("NCM com ID %s não encontrado", id)
	}

	return nil
}

func (r *ncmRepository) GetNCMByCode(code string) (*model.NCM, error) {
	query := `
		SELECT id, code, code_no_symbols, description, initial_date, final_date, 
			   type_year_ini, number_ato_ini, year_ato_ini
		FROM ncm 
		WHERE code = $1
	`

	row := r.db.QueryRow(query, code)

	ncm := &model.NCM{}
	err := row.Scan(&ncm.ID, &ncm.Code, &ncm.CodeNoSimbols, &ncm.Description,
		&ncm.InitialDate, &ncm.FinalDate, &ncm.TypeYearIni, &ncm.NumberAtoIni, &ncm.YearAtoIni)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("NCM com código %s não encontrado", code)
		}
		return nil, fmt.Errorf("erro ao buscar NCM por código: %w", err)
	}

	return ncm, nil
}

func (r *ncmRepository) GetNCMByID(id string) (*model.NCM, error) {
	query := `
		SELECT id, code, code_no_symbols, description, initial_date, final_date, 
			   type_year_ini, number_ato_ini, year_ato_ini
		FROM ncm 
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	ncm := &model.NCM{}
	err := row.Scan(&ncm.ID, &ncm.Code, &ncm.CodeNoSimbols, &ncm.Description,
		&ncm.InitialDate, &ncm.FinalDate, &ncm.TypeYearIni, &ncm.NumberAtoIni, &ncm.YearAtoIni)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("NCM com ID %s não encontrado", id)
		}
		return nil, fmt.Errorf("erro ao buscar NCM por ID: %w", err)
	}

	return ncm, nil
}

func (r *ncmRepository) GetNCMByText(text string) (*model.NCM, error) {
	query := `
		SELECT id, code, code_no_symbols, description, initial_date, final_date, 
			   type_year_ini, number_ato_ini, year_ato_ini
		FROM ncm 
		WHERE LOWER(description) LIKE LOWER($1)
		LIMIT 1
	`

	row := r.db.QueryRow(query, "%"+text+"%")

	ncm := &model.NCM{}
	err := row.Scan(&ncm.ID, &ncm.Code, &ncm.CodeNoSimbols, &ncm.Description,
		&ncm.InitialDate, &ncm.FinalDate, &ncm.TypeYearIni, &ncm.NumberAtoIni, &ncm.YearAtoIni)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("NCM com texto '%s' não encontrado", text)
		}
		return nil, fmt.Errorf("erro ao buscar NCM por texto: %w", err)
	}

	return ncm, nil
}

func (r *ncmRepository) GetNCMSByCodes(codes []string) ([]*model.NCM, error) {
	if len(codes) == 0 {
		return []*model.NCM{}, nil
	}

	placeholders := make([]string, len(codes))
	args := make([]interface{}, len(codes))

	for i, code := range codes {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = code
	}

	query := fmt.Sprintf(`
		SELECT id, code, code_no_symbols, description, initial_date, final_date, 
			   type_year_ini, number_ato_ini, year_ato_ini
		FROM ncm 
		WHERE code IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar NCMs por códigos: %w", err)
	}
	defer rows.Close()

	var ncms []*model.NCM
	for rows.Next() {
		ncm := &model.NCM{}
		err := rows.Scan(&ncm.ID, &ncm.Code, &ncm.CodeNoSimbols, &ncm.Description,
			&ncm.InitialDate, &ncm.FinalDate, &ncm.TypeYearIni, &ncm.NumberAtoIni, &ncm.YearAtoIni)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear NCM: %w", err)
		}
		ncms = append(ncms, ncm)
	}

	return ncms, nil
}

func (r *ncmRepository) GetNCMSByText(text string) ([]*model.NCM, error) {
	query := `
		SELECT id, code, code_no_symbols, description, initial_date, final_date, 
			   type_year_ini, number_ato_ini, year_ato_ini
		FROM ncm 
		WHERE LOWER(description) LIKE LOWER($1)
	`

	rows, err := r.db.Query(query, "%"+text+"%")
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar NCMs por texto: %w", err)
	}
	defer rows.Close()

	var ncms []*model.NCM
	for rows.Next() {
		ncm := &model.NCM{}
		err := rows.Scan(&ncm.ID, &ncm.Code, &ncm.CodeNoSimbols, &ncm.Description,
			&ncm.InitialDate, &ncm.FinalDate, &ncm.TypeYearIni, &ncm.NumberAtoIni, &ncm.YearAtoIni)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear NCM: %w", err)
		}
		ncms = append(ncms, ncm)
	}

	return ncms, nil
}

func (r *ncmRepository) GetAllNCMs() ([]*model.NCM, error) {
	query := `
		SELECT id, code, code_no_symbols, description, initial_date, final_date, 
			   type_year_ini, number_ato_ini, year_ato_ini
		FROM ncm
		ORDER BY code
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar todos os NCMs: %w", err)
	}
	defer rows.Close()

	var ncms []*model.NCM
	for rows.Next() {
		ncm := &model.NCM{}
		err := rows.Scan(&ncm.ID, &ncm.Code, &ncm.CodeNoSimbols, &ncm.Description,
			&ncm.InitialDate, &ncm.FinalDate, &ncm.TypeYearIni, &ncm.NumberAtoIni, &ncm.YearAtoIni)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear NCM: %w", err)
		}
		ncms = append(ncms, ncm)
	}

	return ncms, nil
}

func (r *ncmRepository) BulkInsertNCMs(ncms []*model.NCM) error {
	if len(ncms) == 0 {
		return nil
	}

	query := `
		INSERT INTO ncm (id, code, code_no_symbols, description, initial_date, final_date, type_year_ini, number_ato_ini, year_ato_ini)
		VALUES %s
		ON CONFLICT (id) DO NOTHING
	`

	values := make([]string, len(ncms))
	args := make([]interface{}, 0, len(ncms)*9)

	for i, ncm := range ncms {
		values[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9)
		args = append(args, uuid.New().String(), ncm.Code, ncm.CodeNoSimbols, ncm.Description,
			ncm.InitialDate, ncm.FinalDate, ncm.TypeYearIni, ncm.NumberAtoIni, ncm.YearAtoIni)
	}
	query = fmt.Sprintf(query, strings.Join(values, ", "))
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("erro ao inserir NCMs em massa: %w", err)
	}
	return nil
}
