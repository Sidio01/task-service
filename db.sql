-- Table: public.tasks
-- DROP TABLE public.tasks;
CREATE TABLE IF NOT EXISTS public.tasks
(
    uuid uuid PRIMARY KEY NOT NULL,
    name character(50),
    text character(100),
    login character(50),
    status character(20) NOT NULL DEFAULT 'created'::bpchar
);

-- Table: public.approvals
-- DROP TABLE public.approvals;
CREATE TABLE IF NOT EXISTS public.approvals
(
    id serial PRIMARY KEY NOT NULL,
    task_uuid uuid NOT NULL,
    approval_login character(50) NOT NULL,
    approved boolean DEFAULT NULL,
    sent boolean DEFAULT NULL,
    n integer,
    CONSTRAINT task_uuid FOREIGN KEY (task_uuid)
        REFERENCES public.tasks (uuid) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
);
