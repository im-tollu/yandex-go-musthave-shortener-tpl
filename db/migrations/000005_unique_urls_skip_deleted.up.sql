alter table URLS
    drop constraint UNQ_URLS_ORIGINAL_URL;


create unique index UNQ_URLS_ORIGINAL_URL on URLS (URLS_ORIGINAL_URL) where (URLS_DELETED = false);