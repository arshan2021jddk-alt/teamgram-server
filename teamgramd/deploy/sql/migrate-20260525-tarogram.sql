-- Rename default system account branding from Teamgram to Tarogram.
-- Safe to run multiple times.

UPDATE users
SET first_name = 'Tarogram',
    username = 'tarogram',
    updated_at = NOW()
WHERE id = 777000
  AND (first_name <> 'Tarogram' OR username <> 'tarogram');
