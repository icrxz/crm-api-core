ALTER TABLE IF EXISTS cases
  ADD CONSTRAINT cases_reference_unique_constraint UNIQUE (external_reference, contractor_id);