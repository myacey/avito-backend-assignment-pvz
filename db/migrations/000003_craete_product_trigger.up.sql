CREATE OR REPLACE FUNCTION check_reception_status_before_product_insert()
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$$
DECLARE
    rec_status TEXT;
BEGIN
    SELECT status INTO rec_status FROM receptions WHERE id = NEW.reception_id;

    IF rec_status IS DISTINCT FROM 'in_progress' THEN
        RAISE EXCEPTION 'cannot add product to finished reception';
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trigger_check_reception_status
    BEFORE INSERT OR UPDATE
    ON products
    FOR EACH ROW
    EXECUTE PROCEDURE check_reception_status_before_product_insert();
