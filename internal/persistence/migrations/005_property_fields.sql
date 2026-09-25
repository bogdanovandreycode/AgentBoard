ALTER TABLE property_definitions ADD COLUMN placeholder TEXT NOT NULL DEFAULT '';
ALTER TABLE property_definitions ADD COLUMN regex TEXT NOT NULL DEFAULT '';
ALTER TABLE property_definitions ADD COLUMN default_value TEXT NOT NULL DEFAULT '';
CREATE TRIGGER IF NOT EXISTS task_property_defaults AFTER INSERT ON tasks
BEGIN
 INSERT INTO task_property_values(task_id,property_definition_id,value,updated_at)
 SELECT NEW.id,id,default_value,NEW.created_at FROM property_definitions
 WHERE project_id=NEW.project_id AND default_value!='';
END;
