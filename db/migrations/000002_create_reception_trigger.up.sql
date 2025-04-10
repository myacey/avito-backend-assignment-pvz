CREATE OR REPLACE FUNCTION add_reception_check()
    RETURNS TRIGGER
    LANGUAGE plpgsql
    AS
$$
BEGIN
    IF (SELECT COUNT(*) FROM receptions
        WHERE pvz_id = NEW.pvz_id AND status = 'in_progress' AND id != NEW.id) > 0 THEN
        RAISE EXCEPTION 'cannot create new reception while other in progress';
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trigger_add_reception
    BEFORE INSERT OR UPDATE
    ON receptions
    FOR EACH ROW
    EXECUTE PROCEDURE add_reception_check();