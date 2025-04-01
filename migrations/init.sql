DROP TABLE IF EXISTS flights;

CREATE TABLE flights (
                         id SERIAL PRIMARY KEY,
                         destination_from VARCHAR(100) NOT NULL,
                         destination_to VARCHAR(100) NOT NULL,
                         delete_at BIGINT DEFAULT 0
);
INSERT INTO flights (destination_from, destination_to, delete_at) VALUES
                                                                      ('Москва', 'Париж', 0),
                                                                      ('Берлин', 'Лондон', 0),
                                                                      ('Нью-Йорк', 'Токио', 0),
                                                                      ('Сидней', 'Мельбурн', 0),
                                                                      ('Рим', 'Афины', 0),
                                                                      ('Пекин', 'Шанхай', 0),
                                                                      ('Кейптаун', 'Йоханнесбург', 0),
                                                                      ('Мехико', 'Лима', 0),
                                                                      ('Торонто', 'Монреаль', 0),
                                                                      ('Дубай', 'Дели', 0);
