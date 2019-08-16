INSERT INTO character (name)
VALUES ('basic'),
       ('fox'),
       ('link');

INSERT INTO lesson (name, character_id, number, gif, description, learning_time_seconds, training_time_seconds, test_time_seconds)
VALUES (
            'Short Hop', 
            4, 
            1, 
            'https://ftp.crabbymonkey.org/smash/smash_gifs/smash_examples/example_short_hop.gif', 
            'This is where I type explenations, Lorem ipsum dolor sit amet, fabulas nusquam facilisi per cu, ex ius voluptua principes. Quo te simul nullam. Illud aperiam accusamus mel no. Ex oporteat perfecto petentium qui, meis solum utamur sit te, per reque eligendi appellantur ei. Posse dictas laoreet pri ut, vide tamquam quaeque at his. Eu his bonorum dolorum, est vidisse discere verterem cu. Vim an veritus adipisci. An quaeque alienum electram vis, possim diceret efficiendi ex vis. Id offendit moderatius intellegam pro, ne usu atqui verterem philosophia, sit eu feugiat gloriatur expetendis. Vix ei aperiri scripserit.',
            300,
            1800,
            60
        ),
       (
            'Short Hop, Fast Fall', 
            4, 
            2, 
            'https://ftp.crabbymonkey.org/smash/smash_gifs/smash_examples/example_short_hop.gif', 
            'This is where I type explenations Again! Lorem ipsum dolor sit amet, fabulas nusquam facilisi per cu, ex ius voluptua principes. Quo te simul nullam. Illud aperiam accusamus mel no. Ex oporteat perfecto petentium qui, meis solum utamur sit te, per reque eligendi appellantur ei. Posse dictas laoreet pri ut, vide tamquam quaeque at his. Eu his bonorum dolorum, est vidisse discere verterem cu. Vim an veritus adipisci. An quaeque alienum electram vis, possim diceret efficiendi ex vis. Id offendit moderatius intellegam pro, ne usu atqui verterem philosophia, sit eu feugiat gloriatur expetendis. Vix ei aperiri scripserit.',
            300,
            1800,
            60
        ); 