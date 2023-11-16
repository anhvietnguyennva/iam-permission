CREATE TABLE services
(
    id          varchar(50) PRIMARY KEY,
    namespace   varchar(255) NOT NULL, -- does not allow to be updated
    description varchar(255),
    created_by  varchar(255),
    updated_by  varchar(255),
    created_at  bigint,
    updated_at  bigint
);
CREATE UNIQUE INDEX services_unique_namespace ON services (LOWER(namespace));
CREATE INDEX services_namespace ON services (namespace);

CREATE TABLE relation_definitions
(
    id          varchar(50) PRIMARY KEY,
    service_id  varchar(50) REFERENCES services (id) ON DELETE CASCADE, -- does not allow to be updated
    namespace   varchar(255) NOT NULL,                                  -- does not allow to be updated
    relation    varchar(255) NOT NULL,                                  -- does not allow to be updated
    description varchar(255),
    created_by  varchar(255),
    updated_by  varchar(255),
    created_at  bigint,
    updated_at  bigint
);
CREATE UNIQUE INDEX relation_definitions_unique_service_id_relation ON relation_definitions (service_id, LOWER(relation));
CREATE INDEX relation_definitions_namespace_relation ON relation_definitions (namespace, relation);

CREATE TABLE relation_configurations
(
    id                            varchar(50) PRIMARY KEY,
    service_id                    varchar(50) REFERENCES services (id) ON DELETE CASCADE, -- does not allow to be updated
    namespace                     varchar(255) NOT NULL,                                  -- does not allow to be updated
    parent_relation_definition_id varchar(255) REFERENCES relation_definitions (id) ON DELETE CASCADE,
    parent_relation               varchar(255) NOT NULL,
    child_relation_definition_id  varchar(255) REFERENCES relation_definitions (id) ON DELETE CASCADE,
    child_relation                varchar(255) NOT NULL,
    created_by                    varchar(255),
    updated_by                    varchar(255),
    created_at                    bigint,
    updated_at                    bigint,
    UNIQUE (service_id, parent_relation_definition_id, child_relation_definition_id)
);

CREATE TABLE subject_relation_tuples
(
    id                     varchar(50) PRIMARY KEY,
    service_id             varchar(50) REFERENCES services (id) ON DELETE CASCADE, -- does not allow to be updated
    relation_definition_id varchar(50) REFERENCES relation_definitions (id) ON DELETE CASCADE,
    namespace              varchar(255) NOT NULL,                                  -- does not allow to be updated
    object                 varchar(255) NOT NULL,
    relation               varchar(255) NOT NULL,
    subject_id             varchar(255) NOT NULL,                                  -- IAM Identity ID or IAM Client ID in case of Client Credentials
    created_by             varchar(255),
    updated_by             varchar(255),
    created_at             bigint,
    updated_at             bigint
);
CREATE UNIQUE INDEX subject_relation_tuples_unique_namespace_object_relation_subject_id ON subject_relation_tuples (LOWER(namespace), LOWER(object), LOWER(relation), subject_id);

CREATE TABLE subject_sets
(
    id                     varchar(50) PRIMARY KEY,
    service_id             varchar(50) REFERENCES services (id) ON DELETE CASCADE,
    relation_definition_id varchar(50) REFERENCES relation_definitions (id) ON DELETE CASCADE,
    namespace              varchar(255) NOT NULL,
    object                 varchar(255) NOT NULL,
    relation               varchar(255) NOT NULL,
    description            varchar(500),
    created_by             varchar(255),
    updated_by             varchar(255),
    created_at             bigint,
    updated_at             bigint
);
CREATE UNIQUE INDEX subject_sets_unique_namespace_object_relation ON subject_sets (LOWER(namespace), LOWER(object), LOWER(relation));
CREATE INDEX subject_sets_namespace_object_relation ON subject_sets (namespace, object, relation);

CREATE TABLE subject_set_relation_tuples
(
    id                     varchar(50) PRIMARY KEY,
    service_id             varchar(50) REFERENCES services (id) ON DELETE CASCADE,
    subject_set_id         varchar(50) REFERENCES subject_sets (id) ON DELETE CASCADE,
    relation_definition_id varchar(50) REFERENCES relation_definitions (id) ON DELETE CASCADE,
    namespace              varchar(255) NOT NULL,
    object                 varchar(255) NOT NULL,
    relation               varchar(255) NOT NULL,
    subject_set_namespace  varchar(255) NOT NULL,
    subject_set_object     varchar(255) NOT NULL,
    subject_set_relation   varchar(255) NOT NULL,
    created_by             varchar(255),
    updated_by             varchar(255),
    created_at             bigint,
    updated_at             bigint
);
CREATE UNIQUE INDEX subject_set_relation_tuples_unique_namespace_object_relation_subject_set_namespace_subject_set_object_subject_set_relation ON subject_set_relation_tuples (LOWER(namespace), LOWER(object), LOWER(relation), LOWER(subject_set_namespace), LOWER(subject_set_object), LOWER(subject_set_relation));
