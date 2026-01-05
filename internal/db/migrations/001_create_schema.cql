CREATE TYPE IF NOT EXISTS check_config (
    timeout int,
    method text,
    headers map<text, text>,
    body text,
    follow_redirects boolean,
    verify_ssl boolean
);

CREATE TYPE IF NOT EXISTS agent_capabilities (
    check_types set<text>,
    regions set<text>,
    max_concurrent_checks int
);

CREATE TABLE IF NOT EXISTS targets (
    id uuid PRIMARY KEY,
    user_id uuid,
    name text,
    type text,  -- http, tcp, ping, dns
    url text,
    config frozen<check_config>,
    interval_seconds int,
    locations set<text>,  -- geographic regions
    tags set<text>,
    enabled boolean,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS targets_by_user (
    user_id uuid,
    target_id uuid,
    name text,
    type text,
    enabled boolean,
    created_at timestamp,
    PRIMARY KEY ((user_id), target_id)
);

CREATE TABLE IF NOT EXISTS agents (
    id uuid PRIMARY KEY,
    name text,
    version text,
    location text,  -- region/country
    ip inet,
    public_key text,
    status text,  -- active, inactive, suspended, validating
    capabilities frozen<agent_capabilities>,
    last_heartbeat timestamp,
    registered_at timestamp,
    tags set<text>
);

CREATE TABLE IF NOT EXISTS agents_by_status (
    status text,  -- active, inactive, suspended, validating
    agent_id uuid,
    location text,
    last_heartbeat timestamp,
    PRIMARY KEY ((status), agent_id)
);

CREATE TABLE IF NOT EXISTS check_results (
    target_id uuid,
    bucket date,  -- daily bucket for partitioning
    check_time timestamp,
    agent_id uuid,
    status text,  -- up, down, timeout, error
    response_time_ms int,
    response_code int,
    error_message text,
    ssl_expiry timestamp,
    PRIMARY KEY ((target_id, bucket), check_time, agent_id)
) WITH CLUSTERING ORDER BY (check_time DESC, agent_id ASC)
  AND default_time_to_live = 2592000;

CREATE TABLE IF NOT EXISTS latest_check_results (
    target_id uuid,
    check_time timestamp,
    agent_id uuid,
    status text,  -- up, down, timeout, error
    response_time_ms int,
    response_code int,
    PRIMARY KEY ((target_id), check_time, agent_id)
) WITH CLUSTERING ORDER BY (check_time DESC);

CREATE TABLE IF NOT EXISTS incidents (
    target_id uuid,
    started_at timestamp,
    incident_id uuid,
    ended_at timestamp,
    status text,  -- active, resolved, acknowledged
    title text,
    description text,
    acknowledged_at timestamp,
    acknowledged_by uuid,
    resolved_at timestamp,
    resolved_by uuid,
    downtime_seconds int,
    affected_checks int,
    PRIMARY KEY ((target_id), started_at, incident_id)
) WITH CLUSTERING ORDER BY (started_at DESC);

CREATE TABLE IF NOT EXISTS active_incidents (
    target_id uuid,
    incident_id uuid,
    started_at timestamp,
    title text,
    PRIMARY KEY ((target_id), incident_id)
);

CREATE TABLE IF NOT EXISTS status_history (
    target_id uuid,
    bucket date,
    changed_at timestamp,
    old_status text,
    new_status text,
    reason text,
    agent_id uuid,
    PRIMARY KEY ((target_id, bucket), changed_at)
) WITH CLUSTERING ORDER BY (changed_at DESC);
