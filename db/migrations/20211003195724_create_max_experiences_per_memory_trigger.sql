-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION max_experiences_per_memory ()
    RETURNS TRIGGER
    AS $$
DECLARE
    max_experience_count integer := 3;
    experience_count integer := 0;
    must_check boolean := FALSE;
BEGIN
    IF TG_OP = 'INSERT' THEN
        must_check := TRUE;
    END IF;
    IF TG_OP = 'UPDATE' THEN
        IF (NEW.memory_id != OLD.memory_id) THEN
            must_check := TRUE;
        END IF;
    END IF;
    IF must_check THEN
        LOCK TABLE experiences IN exclusive mode;
        SELECT
            INTO experience_count count(*)
        FROM
            experiences
        WHERE
            memory_id = NEW.memory_id;
        IF experience_count >= max_experience_count THEN
            RAISE EXCEPTION
                USING ERRCODE = 'TH001', MESSAGE = format('cannot insert more than %s experiences per memory', max_experience_count);
            END IF;
        END IF;
        RETURN new;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER max_experiences_per_memory
    BEFORE INSERT OR UPDATE ON experiences
    FOR EACH ROW
    EXECUTE PROCEDURE max_experiences_per_memory ();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER max_experiences_per_memory ON experiences;

DROP FUNCTION max_experiences_per_memory ();

-- +goose StatementEnd
