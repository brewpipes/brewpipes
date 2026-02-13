CREATE INDEX IF NOT EXISTS batch_process_phase_batch_id_phase_at_idx
  ON batch_process_phase(batch_id, phase_at DESC)
  WHERE deleted_at IS NULL;
