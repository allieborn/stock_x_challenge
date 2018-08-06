This is a simple web app, written in Go and using mysql, which allows users to rate how 'true' a shoe fits.

To get this project running locally on your machine, you must have mysql installed. In mysql, run the following commands:

         drop database if exists shoes;
         create database shoes;

         use shoes;

         drop table if exists shoe_catalog;
         create table shoe_catalog (id varchar(10), shoeName varchar(50), designer varchar(50), price varchar(50));

         insert into shoe_catalog (id, shoeName, designer, price) values ('101', 'Pigalle Follies', 'Christian Louboutin', '765');
         insert into shoe_catalog (id, shoeName, designer, price) values ('102', 'Point Toe Ballet Flats', 'Dolce & Gabana', '140');


        drop table if exists true_to_size_scoring;
        create table true_to_size_scoring (shoeID varchar(10), score varchar(10), user varchar(50), date varchar(50));

        insert into true_to_size_scoring (shoeID, score,  date) values ('101', '5','05-19-2018');
        insert into true_to_size_scoring (shoeID, score,  date) values ('102', '1','05-19-2018');

Next, in the main.go file, where the variable 'connStr' is initialized, replace the placeholder values with your sql
username and your sql port address.

Run the project, and visit localhost:8080 in your browser. You will see the web app there.
