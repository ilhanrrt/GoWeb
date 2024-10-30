-- phpMyAdmin SQL Dump
-- version 5.1.2
-- https://www.phpmyadmin.net/
--
-- Anamakine: localhost:3306
-- Üretim Zamanı: 16 Ağu 2024, 08:52:29
-- Sunucu sürümü: 5.7.24
-- PHP Sürümü: 8.3.1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Veritabanı: `blog`
--

-- --------------------------------------------------------

--
-- Tablo için tablo yapısı `categories`
--

CREATE TABLE `categories` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` longtext,
  `slug` longtext
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Tablo döküm verisi `categories`
--

INSERT INTO `categories` (`id`, `created_at`, `updated_at`, `deleted_at`, `title`, `slug`) VALUES
(2, '2024-08-05 23:48:49.435', '2024-08-05 23:48:49.435', '2024-08-05 23:52:19.101', 'Teknoloji', 'teknoloji'),
(3, '2024-08-05 23:48:59.419', '2024-08-05 23:48:59.419', '2024-08-05 23:52:15.761', 'Yazılım', 'yazilim'),
(4, '2024-08-05 23:52:23.737', '2024-08-05 23:52:23.737', NULL, 'Teknoloji', 'teknoloji'),
(5, '2024-08-05 23:52:25.385', '2024-08-05 23:52:25.385', NULL, 'Yazılım', 'yazilim'),
(6, '2024-08-05 23:52:30.997', '2024-08-05 23:52:30.997', NULL, 'Veri Bilimi', 'veri-bilimi'),
(7, '2024-08-05 23:52:35.884', '2024-08-05 23:52:35.884', NULL, 'Donanım', 'donanim');

-- --------------------------------------------------------

--
-- Tablo için tablo yapısı `posts`
--

CREATE TABLE `posts` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` longtext,
  `slug` longtext,
  `description` longtext,
  `content` longtext,
  `picture_url` longtext,
  `category_id` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Tablo döküm verisi `posts`
--

INSERT INTO `posts` (`id`, `created_at`, `updated_at`, `deleted_at`, `title`, `slug`, `description`, `content`, `picture_url`, `category_id`) VALUES
(1, '2024-07-30 15:36:06.689', '2024-07-30 16:17:09.107', '2024-08-07 16:06:56.000', 'Phyton ile web programlama', 'web', 'test', '', '', 0),
(2, '2024-07-31 17:17:06.825', '2024-08-06 16:18:32.481', NULL, 'go ile web programlama', 'go-ile-web-programlama', 'go ile web programlama', '<ul><li><font color=\"#000000\" style=\"background-color: rgb(255, 255, 0);\">go ile web programlama</font><br></li></ul>', 'uploads/2018-z-performance-bmw-m5-g30-04.jpg', 4),
(3, '2024-08-02 11:33:34.000', '2024-08-02 11:33:34.000', NULL, 'asdfasdf', 'asdf', 'asdf', 'asdf', NULL, NULL),
(70, '2024-08-03 11:34:41.684', '2024-08-03 11:34:41.684', NULL, 'adsf', 'adsf', 'aasdf', '<p>asdf</p>', 'uploads/Ekran görüntüsü 2024-08-02 212842.png', 1),
(71, '2024-08-03 11:48:48.675', '2024-08-06 13:12:00.592', NULL, 'asdf', 'asdf', 'asdf', '<p>asdf</p>', 'uploads/Ekran görüntüsü 2024-08-02 212842.png', 6),
(72, '2024-08-03 11:49:30.592', '2024-08-06 16:07:42.120', NULL, 'asdf', 'asdf', 'asdf', '<p>asdf</p>', 'uploads/Ekran görüntüsü 2024-08-02 212842.png', 5),
(73, '2024-08-03 14:59:25.546', '2024-08-16 11:29:44.310', NULL, 'kayıt düzenleme', 'kayit-duzenleme', 'asd', '<p>asd</p>', 'uploads/Ekran görüntüsü 2024-08-02 212842.png', 7),
(74, '2024-08-06 13:13:27.689', '2024-08-06 16:08:54.796', NULL, 'veri', 'veri', 'veri bilimi', '<p>gfhjnöç</p>', 'uploads/512416.jpg', 4),
(75, '2024-08-06 16:19:27.059', '2024-08-06 16:20:00.458', NULL, 'go ile web programlama', 'go-ile-web-programlama', 'go ile web programlama', '<h1><u><b><font color=\"#000000\" style=\"background-color: rgb(255, 255, 0);\"><span style=\"font-family: &quot;Arial Black&quot;;\">go ile web programlama</span></font></b></u></h1>', 'uploads/2018-z-performance-bmw-m5-g30-04.jpg', 6),
(76, '2024-08-15 11:12:00.792', '2024-08-15 11:12:00.792', NULL, 'go ile filtre', 'go-ile-filtre', 'Filtreleme öğreniyoruz', '<p><font color=\"#000000\" style=\"background-color: rgb(255, 255, 0);\"><u><b>go ile filtre öğreniyoruz</b></u></font><br></p>', 'uploads/512416.jpg', 5),
(77, '2024-08-15 14:55:06.634', '2024-08-15 14:55:06.634', NULL, 'go ile moderatör denemesi', 'go-ile-moderator-denemesi', 'go ile moderatör denemesi', '<p>go ile moderatör denemesi<br></p>', 'uploads/413455.jpg', 4),
(78, '2024-08-15 14:58:18.680', '2024-08-15 14:58:18.680', NULL, 'go ile moderatör denemesi 2', 'go-ile-moderator-denemesi-2', 'go ile moderatör denemesi 2', '<div style=\"background-color: rgb(30, 31, 34);\"><pre style=\"\"><font color=\"#c77dbb\" face=\"JetBrains Mono, monospace\"><i>go ile moderatör denemesi 2</i></font><font color=\"#bcbec4\" face=\"JetBrains Mono, monospace\"><br></font></pre></div>', 'uploads/926807.jpg', 5),
(79, '2024-08-15 14:58:49.126', '2024-08-15 14:58:49.126', NULL, 'go ile moderatör denemesi 2', 'go-ile-moderator-denemesi-2', 'go ile moderatör denemesi 2', '<p>go ile moderatör denemesi 2<br></p>', 'uploads/512422.jpg', 7),
(80, '2024-08-15 15:00:54.035', '2024-08-15 15:00:54.035', '2024-08-15 19:44:00.745', 'go ile moderatör denemesi 2', 'go-ile-moderator-denemesi-2', 'go ile moderatör denemesi 2', '<p>go ile moderatör denemesi 2<br></p>', 'uploads/bmw1310.jpg', 5),
(81, '2024-08-15 19:07:16.299', '2024-08-15 19:07:16.299', '2024-08-15 19:40:30.194', 'go ile moderatör denemesi 2', 'go-ile-moderator-denemesi-2', 'go ile moderatör denemesi 2', '<p>go ile moderatör denemesi 2<br></p>', 'uploads/b7 alpina blue.jpg', 6),
(82, '2024-08-15 19:39:25.372', '2024-08-15 19:39:25.372', NULL, 'go ile moderatör denemesi bilmem kaç', 'go-ile-moderator-denemesi-bilmem-kac', 'go ile moderatör denemesi bilmem kaç', '<p>go ile moderatör denemesi bilmem kaç<br></p>', 'uploads/512416.jpg', 6),
(83, '2024-08-15 19:44:16.359', '2024-08-15 19:44:16.359', '2024-08-15 19:51:13.583', 'go ile moderatör denemesi bilmem kaç', 'go-ile-moderator-denemesi-bilmem-kac', 'go ile moderatör denemesi bilmem kaç', '<p>go ile moderatör denemesi bilmem kaç<br></p>', 'uploads/512416.jpg', 5);

-- --------------------------------------------------------

--
-- Tablo için tablo yapısı `users`
--

CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` longtext,
  `password` longtext,
  `user_type` longtext
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Tablo döküm verisi `users`
--

INSERT INTO `users` (`id`, `created_at`, `updated_at`, `deleted_at`, `username`, `password`, `user_type`) VALUES
(1, NULL, NULL, NULL, 'admin', '86362fda56389b5a3613602561362b6c793ea875ef19bd96b017a11867436597', 'admin'),
(2, NULL, NULL, NULL, 'moderatör', '86362fda56389b5a3613602561362b6c793ea875ef19bd96b017a11867436597', 'viewer'),
(3, NULL, NULL, NULL, 'editör', '86362fda56389b5a3613602561362b6c793ea875ef19bd96b017a11867436597', 'editor');

--
-- Dökümü yapılmış tablolar için indeksler
--

--
-- Tablo için indeksler `categories`
--
ALTER TABLE `categories`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_categories_deleted_at` (`deleted_at`);

--
-- Tablo için indeksler `posts`
--
ALTER TABLE `posts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_posts_deleted_at` (`deleted_at`);

--
-- Tablo için indeksler `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_users_deleted_at` (`deleted_at`);

--
-- Dökümü yapılmış tablolar için AUTO_INCREMENT değeri
--

--
-- Tablo için AUTO_INCREMENT değeri `categories`
--
ALTER TABLE `categories`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- Tablo için AUTO_INCREMENT değeri `posts`
--
ALTER TABLE `posts`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=84;

--
-- Tablo için AUTO_INCREMENT değeri `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
