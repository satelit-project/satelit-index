-- +goose Up

create extension if not exists "uuid-ossp";

-- +goose StatementBegin
create or replace function manage_updated_at(_tbl regclass) returns void as $$
begin
    execute format('create trigger set_updated_at before update on %s
                    for each row execute procedure set_updated_at()', _tbl);
end;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
create or replace function unmanage_updated_at(_tbl regclass) returns void as $$
begin
    execute format('drop trigger set_updated_at on %s', _tbl);
end;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
create or replace function set_updated_at() returns trigger as $$
begin
    if (
        new is distinct from old and
        new.updated_at is not distinct from old.updated_at
    ) then
        new.updated_at := current_timestamp;
    end if;
    return new;
end;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose Down

drop extension if exists "uuid-ossp";
drop function if exists manage_updated_at;
drop function if exists unmanage_updated_at;
drop function if exists set_updated_at;
