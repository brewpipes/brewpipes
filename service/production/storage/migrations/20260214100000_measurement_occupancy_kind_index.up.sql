CREATE INDEX idx_measurement_occupancy_kind_observed
  ON measurement(occupancy_id, kind, observed_at DESC)
  WHERE deleted_at IS NULL;
