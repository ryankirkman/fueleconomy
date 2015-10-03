-- +migrate Up
CREATE TABLE vehicles (
    id                       serial primary key,
    updated                  timestamptz default now(),
    atv_type                 varchar(255),
    charge_time_120v         float8,
    charge_time_240v         float8,
    charge_time_240vb        float8,
    charger_240v_dscr        varchar(255),
    charger_240vb_dscr       varchar(255),
    cylinders                integer,
    drive_axle_type          varchar(255),
    e_city                   float8,
    e_comb                   float8,
    e_highway                float8,
    e_motor                  varchar(255),
    eng_displacement         float8,
    eng_dscr                 varchar(255),
    eng_id                   integer,
    epa_created_on           timestamptz,
    epa_id                   integer unique,
    epa_modified_on          timestamptz,
    f1_barrels_per_year      float8,
    f1_co2                   float8,
    f1_co2_tailpipe          float8,
    f1_fuel_cost             integer,
    f1_fuel_type             varchar(255),
    f1_ghg_score             integer,
    f1_mpg_city              float8,
    f1_mpg_city_unadj        float8,
    f1_mpg_city_unrounded    float8,
    f1_mpg_comb              float8,
    f1_mpg_comb_unrounded    float8,
    f1_mpg_highway           float8,
    f1_mpg_highway_unrounded float8,
    f1_mpg_highway_unadj     float8,
    f1_range                 float8,
    f2_barrels_per_year      float8,
    f2_co2                   float8,
    f2_co2_tailpipe          float8,
    f2_fuel_cost             integer,
    f2_fuel_type             varchar(255),
    f2_ghg_score             integer,
    f2_mpg_city              float8,
    f2_mpg_city_unadj        float8,
    f2_mpg_city_unrounded    float8,
    f2_mpg_comb              float8,
    f2_mpg_comb_unrounded    float8,
    f2_mpg_highway           float8,
    f2_mpg_highway_unrounded float8,
    f2_mpg_highway_unadj     float8,
    f2_range                 float8,
    f2_range_city            float8,
    f2_range_highway         float8,
    fuel_economy_score       float8,
    fuel_type                varchar(255),
    is_guzzler               boolean,
    is_phev_blended          boolean,
    has_mpg_data             boolean,
    has_supercharger         boolean,
    has_turbocharger         boolean,
    luggage_volume_2door     integer,
    luggage_volume_4door     integer,
    luggage_volume_hatch     integer,
    make                     varchar(255),
    manufacturer_code        varchar(255),
    model                    varchar(255),
    mpg_data                 varchar(255),
    passenger_volume_2door   integer,
    passenger_volume_4door   integer,
    passenger_volume_hatch   integer,
    phev_cd_city             float8,
    phev_cd_comb             float8,
    phev_cd_highway          float8,
    phev_mpg_city            float8,
    phev_mpg_comb            float8,
    phev_mpg_highway         float8,
    phev_uf_city             float8,
    phev_uf_comb             float8,
    phev_uf_highway          float8,
    size_class               varchar(255),
    start_stop               varchar(255),
    trans_dscr               varchar(255),
    transition               varchar(255),
    year                     integer,
    you_save_spend           integer
);

GRANT SELECT, UPDATE, INSERT, DELETE ON vehicles TO api;
GRANT USAGE, SELECT, UPDATE ON vehicles_id_seq TO api;

CREATE TABLE emissions_info (
    id                       serial primary key,
    updated                  timestamptz default now(),
    emission_std_code        varchar(255),
    emission_std_txt         varchar(255),
    engine_family_id         varchar(255),
    epa_id                   integer references vehicles(epa_id) on delete cascade on update cascade,
    f1_smog_rating           integer,
    f2_smog_rating           integer,
    sales_area               integer,
    smartway_score           integer
);

GRANT SELECT, UPDATE, INSERT, DELETE ON emissions_info TO api;
GRANT USAGE, SELECT, UPDATE ON emissions_info_id_seq TO api;

CREATE TABLE fuel_prices (
    id                       serial primary key,
    updated                  timestamptz default now(),
    cng                      float8,
    diesel                   float8,
    e85                      float8,
    electricity              float8,
    gas_midgrade             float8,
    gas_premium              float8,
    gas_regular              float8,
    liquid_propane           float8
);

GRANT SELECT, UPDATE, INSERT, DELETE ON fuel_prices TO api;
GRANT USAGE, SELECT, UPDATE ON fuel_prices_id_seq TO api;

-- +migrate StatementBegin
CREATE FUNCTION upsert(update_sql TEXT, insert_sql TEXT) RETURNS VOID AS
$$
BEGIN
    LOOP
        EXECUTE update_sql;
        IF found THEN
            RETURN;
        END IF;
        BEGIN
            EXECUTE insert_sql;
            RETURN;
        EXCEPTION WHEN unique_violation THEN
            -- do nothing, and loop to try the UPDATE again
        END;
    END LOOP;
END;
$$
LANGUAGE plpgsql;
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE emissions_info;
DROP TABLE vehicles;
DROP TABLE fuel_prices;
DROP FUNCTION upsert(update_sql TEXT, insert_sql TEXT);
