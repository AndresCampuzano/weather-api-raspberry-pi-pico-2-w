-- Reset triggers for weather table

-- Drop existing triggers
DROP TRIGGER IF EXISTS weather_updated_at_trigger ON weather;
DROP TRIGGER IF EXISTS weather_created_at_trigger ON weather;

-- Drop existing functions
DROP FUNCTION IF EXISTS update_timestamp();
DROP FUNCTION IF EXISTS set_created_at();

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    -- Only set updated_at when the record is actually modified
    IF OLD.temperature <> NEW.temperature OR OLD.humidity <> NEW.humidity THEN
        NEW.updated_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create updated_at trigger
CREATE TRIGGER weather_updated_at_trigger
BEFORE UPDATE ON weather
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Create created_at trigger function
CREATE OR REPLACE FUNCTION set_created_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create created_at trigger
CREATE TRIGGER weather_created_at_trigger
BEFORE INSERT ON weather
FOR EACH ROW
EXECUTE FUNCTION set_created_at();

-- Update all existing records to set updated_at to NULL if it equals created_at
UPDATE weather SET updated_at = NULL WHERE updated_at = created_at;