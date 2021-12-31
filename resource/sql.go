package resource

const  CreateTableSqls  = `
create table t_idc (id uuid primary key, create_at int, update_at int, description text, name text unique);

create table t_group (id uuid primary key, create_at int, update_at int, description text, name text unique, idc_id uuid, healthy_check_mode text, healthy_check_config text, healthy_num int, un_healthy_num int, healthy_interval int, unhealthy_interval int, lb_mode text, lb_config text, server_tags text, weight int, ip_alloc_num int, default_servers text, port int, CONSTRAINT idc_id_fk FOREIGN KEY  (idc_id)  REFERENCES t_idc (id));

create table t_server_group_ref (id uuid primary key, create_at int, update_at int, description text, server_ip text, group_id uuid, weight int, CONSTRAINT group_id_fk FOREIGN KEY  (group_id)  REFERENCES t_group (id), constraint uk_server_ip_group_id unique(server_ip,group_id));

create table t_pool (id uuid primary key, create_at int, update_at int, description text, name text unique, lb_mode int, lb_config text, fail_over_config text);

create table t_group_pool_ref (id uuid primary key, create_at int, update_at int, description text, group_id uuid, pool_id uuid, weight int, CONSTRAINT pool_id_fk FOREIGN KEY  (pool_id)  REFERENCES t_pool (id), CONSTRAINT group_id_fk FOREIGN KEY  (group_id)  REFERENCES t_group (id), constraint uk_pool_id_group_id unique(pool_id,group_id));

create table t_service (id uuid primary key, create_at int, update_at int, description text, name text unique, ttl int, pool_id uuid);

create table t_route (id uuid primary key, create_at int, update_at int, description text, name text, service_id uuid, rules text, domains text, priority int, CONSTRAINT service_id_fk FOREIGN KEY  (service_id)  REFERENCES t_service (id));

create table t_region (id uuid primary key, create_at int, update_at int, description text, name text unique, code int unique, idc_affinity text);

create table t_country (id uuid primary key, create_at int, update_at int, description text, code int unique, name text unique);

create table t_province (id uuid primary key, create_at int, update_at int, description text, name text unique, code int unique, region_code int, country_code int, idc_affinity text, CONSTRAINT region_code_fk FOREIGN KEY  (region_code)  REFERENCES t_region (code), CONSTRAINT country_code_fk FOREIGN KEY  (country_code)  REFERENCES t_country (code));

create table t_city (id uuid primary key, create_at int, update_at int, description text, name text unique, code int unique, province_code int, idc_affinity text, CONSTRAINT province_code_fk FOREIGN KEY  (province_code)  REFERENCES t_province (code));

create table t_cluster_event (id bigserial primary key, event text, channel text, data text, at int, expire_at int);
create index index_t_cluster_event_at on t_cluster_event(at);
create index index_t_cluster_event_expire_at on t_cluster_event(expire_at);

create table t_custom_param_enum (id uuid primary key, create_at int, update_at int, description text, param_name text, value text, constraint uk_param_name_value unique(param_name,value));

create table t_user (username text unique, password text, type text, id uuid primary key, create_at int, update_at int, description text);

`

const UpgradeSql = `
alter table t_service add column tags text default '' ;
`
