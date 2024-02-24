
CREATE TABLE public.kab_kotas (
    id smallint,
    nama_kab_kota character varying(32) DEFAULT NULL::character varying,
    provinsi_id smallint,
    latitude character varying(1) DEFAULT NULL::character varying,
    longitude character varying(1) DEFAULT NULL::character varying,
    created_at character varying(1) DEFAULT NULL::character varying,
    updated_at character varying(1) DEFAULT NULL::character varying
);



--
-- Name: _kecamatans; Type: TABLE; Schema: public; Owner: rebasedata
--

CREATE TABLE public.kecamatans (
    id integer,
    nama_kecamatan character varying(31) DEFAULT NULL::character varying,
    kab_kota_id smallint,
    latitude character varying(1) DEFAULT NULL::character varying,
    longitude character varying(1) DEFAULT NULL::character varying,
    created_at character varying(1) DEFAULT NULL::character varying,
    updated_at character varying(1) DEFAULT NULL::character varying
);



--
-- Name: _kelurahan_desas; Type: TABLE; Schema: public; Owner: rebasedata
--

CREATE TABLE public.kelurahan_desas (
    id bigint,
    kecamatan_id integer,
    nama_kelurahan_desa character varying(37) DEFAULT NULL::character varying,
    latitude character varying(1) DEFAULT NULL::character varying,
    longitude character varying(1) DEFAULT NULL::character varying,
    created_at character varying(1) DEFAULT NULL::character varying,
    updated_at character varying(1) DEFAULT NULL::character varying
);


--
-- Name: _provinsis; Type: TABLE; Schema: public; Owner: rebasedata
--

CREATE TABLE public.provinsis (
    id smallint,
    nama_provinsi character varying(26) DEFAULT NULL::character varying,
    latitude character varying(1) DEFAULT NULL::character varying,
    longitude character varying(1) DEFAULT NULL::character varying,
    created_at character varying(1) DEFAULT NULL::character varying,
    updated_at character varying(1) DEFAULT NULL::character varying
);
