PGDMP         /        
        {            tgbotdb    15.2    15.2 
    �           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            �           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false                        0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false                       1262    16398    tgbotdb    DATABASE     {   CREATE DATABASE tgbotdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Russian_Russia.1251';
    DROP DATABASE tgbotdb;
                postgres    false            �            1259    16432    orders    TABLE     X   CREATE TABLE public.orders (
    order_id integer NOT NULL,
    order_list integer[]
);
    DROP TABLE public.orders;
       public         heap    postgres    false            �            1259    16431    orders_order_id_seq    SEQUENCE     �   ALTER TABLE public.orders ALTER COLUMN order_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.orders_order_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            public          postgres    false    217            �            1259    16424    users    TABLE     �   CREATE TABLE public.users (
    user_id integer NOT NULL,
    cart integer[],
    order_history_id integer[],
    tguser_id character varying(256),
    address character varying(128),
    last_inserted_id integer,
    mark integer
);
    DROP TABLE public.users;
       public         heap    postgres    false            �            1259    16423    users_user_id_seq    SEQUENCE     �   ALTER TABLE public.users ALTER COLUMN user_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);
            public          postgres    false    215            m           2606    16438    orders pk_order_id 
   CONSTRAINT     V   ALTER TABLE ONLY public.orders
    ADD CONSTRAINT pk_order_id PRIMARY KEY (order_id);
 <   ALTER TABLE ONLY public.orders DROP CONSTRAINT pk_order_id;
       public            postgres    false    217            k           2606    16430    users pk_user_id 
   CONSTRAINT     S   ALTER TABLE ONLY public.users
    ADD CONSTRAINT pk_user_id PRIMARY KEY (user_id);
 :   ALTER TABLE ONLY public.users DROP CONSTRAINT pk_user_id;
       public            postgres    false    215           