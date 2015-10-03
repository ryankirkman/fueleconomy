-- +migrate Up
CREATE INDEX vehicles_id_idx ON vehicles (epa_id);
CREATE INDEX vehicles_make_model_year_idx ON vehicles (make, model, year);
CREATE INDEX emissions_id_idx ON emissions_info (epa_id);

-- +migrate Down
DROP INDEX vehicles_id_idx;
DROP INDEX vehicles_make_model_year_idx;
DROP INDEX emissions_id_idx;
