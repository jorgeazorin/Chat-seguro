-- phpMyAdmin SQL Dump
-- version 4.4.12
-- http://www.phpmyadmin.net
--
-- Servidor: 127.0.0.1
-- Tiempo de generación: 02-05-2016 a las 23:03:00
-- Versión del servidor: 5.6.25
-- Versión de PHP: 5.6.11

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de datos: `sds`
--

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `chat`
--

CREATE TABLE IF NOT EXISTS `chat` (
  `id` int(11) NOT NULL,
  `nombre` varchar(50) DEFAULT NULL
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `chat`
--

INSERT INTO `chat` (`id`, `nombre`) VALUES
(1, NULL),
(2, NULL),
(3, NULL),
(4, NULL),
(5, 'grupo molon :)'),
(6, ''),
(7, '');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `clavesmensajes`
--

CREATE TABLE IF NOT EXISTS `clavesmensajes` (
  `id` int(11) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `clavesmensajes`
--

INSERT INTO `clavesmensajes` (`id`) VALUES
(1),
(2),
(3);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `clavesusuario`
--

CREATE TABLE IF NOT EXISTS `clavesusuario` (
  `idusuario` int(11) NOT NULL,
  `idclavesmensajes` int(11) NOT NULL,
  `claveusuario` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `clavesusuario`
--

INSERT INTO `clavesusuario` (`idusuario`, `idclavesmensajes`, `claveusuario`) VALUES
(1, 1, 'claveusuario1'),
(1, 2, 'otraclave'),
(15, 1, 'minuevaclavemaria');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `mensaje`
--

CREATE TABLE IF NOT EXISTS `mensaje` (
  `id` int(11) NOT NULL,
  `texto` varchar(1000) NOT NULL,
  `emisor` int(11) NOT NULL,
  `chat` int(11) NOT NULL,
  `clave` int(11) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `mensaje`
--

INSERT INTO `mensaje` (`id`, `texto`, `emisor`, `chat`, `clave`) VALUES
(2, 'Hola que tal?? :)', 1, 5, 1),
(3, 'Pero mira quien está por aqui, che!', 13, 5, 1),
(4, 'jajaja si, al final me instale securitychat! ya podemos hablar sin que nos espien!', 1, 5, 1),
(5, 'Pues ya ves, me siento seguro jejeje', 13, 5, 1),
(6, 'Hola amigo 1', 15, 1, 1),
(7, 'Anda amigo 15, que tal andas?? :)', 1, 1, 1),
(8, 'Hola que tal?? :)', 1, 1, 1),
(9, 'Hola que tal?? :)', 1, 1, 1),
(10, 'Hola amigos', 2, 1, 1),
(11, 'Hola', 2, 1, 1),
(12, 'Hola', 2, 1, 2);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `receptoresmensaje`
--

CREATE TABLE IF NOT EXISTS `receptoresmensaje` (
  `idmensaje` int(11) NOT NULL,
  `idreceptor` int(11) NOT NULL,
  `leido` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `receptoresmensaje`
--

INSERT INTO `receptoresmensaje` (`idmensaje`, `idreceptor`, `leido`) VALUES
(2, 13, 0),
(2, 15, 0),
(3, 1, 1),
(3, 15, 1),
(4, 13, 1),
(4, 15, 1),
(5, 1, 1),
(5, 15, 0),
(6, 1, 1),
(7, 15, 0),
(8, 15, 0),
(9, 15, 0),
(10, 1, 0),
(10, 15, 0),
(11, 1, 0),
(11, 15, 0),
(12, 1, 0),
(12, 15, 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuario`
--

CREATE TABLE IF NOT EXISTS `usuario` (
  `id` int(11) NOT NULL,
  `nombre` varchar(15) NOT NULL,
  `clavepubrsa` varchar(255) NOT NULL,
  `claveprivrsa` varchar(255) NOT NULL,
  `clavelogin` blob NOT NULL,
  `salt` blob NOT NULL,
  `clavecifrado` blob NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `usuario`
--

INSERT INTO `usuario` (`id`, `nombre`, `clavepubrsa`, `claveprivrsa`, `clavelogin`, `salt`, `clavecifrado`) VALUES
(1, 'Pepe', 'clave1', 'clave1priv', 0x70, '', ''),
(2, 'Jorge', 'clavepubrsa2', 'clave2priv', 0x6a, '', ''),
(3, 'Encarna', 'clavepubrsa3', 'clave3priv', 0x636c6176653363696672616461, '', ''),
(4, 'lolo', 'clave4rsa', 'clave4priv', 0x636c6176653463696672616461, '', ''),
(5, 'leila', 'clave5rsa', 'clave5priv', 0x636c6176653563696672616461, '', ''),
(13, 'Lucia', 'clave1', 'clave13priv', 0x636c617665313363696672616461, '', ''),
(15, 'Maria', 'clave15pubrsa', 'clave15privrsa', 0x6d, '', ''),
(16, 'alex', 'clavepubrsa', 'claveprivrsa', 0x636c61766563696672616461, '', ''),
(17, 'Marta', 'clavepubrsa', 'claveprivrsa', 0x636c61766563696672616461, '', ''),
(25, 'Prueba', 'Prueba', 'Prueba', 0x8d589a6b39e18b69ff8d02cd239d2716265feed354ddbebe962da0251cf8c70ebb06ad7c2279ff5882e96db5483e3a497eda81b26de37d1a1dd44c02375366fc, 0xde9b8b325ee0a76887a235cba2a30bcfde8d29cef7425fa7b051881e96b6dfb1, 0x1e8545c0bbb13c8b514e006045be320f);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuarioschat`
--

CREATE TABLE IF NOT EXISTS `usuarioschat` (
  `idusuario` int(11) NOT NULL,
  `idchat` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `usuarioschat`
--

INSERT INTO `usuarioschat` (`idusuario`, `idchat`) VALUES
(1, 1),
(15, 1),
(1, 2),
(13, 2),
(1, 5),
(13, 5),
(15, 5),
(1, 6),
(2, 6),
(3, 6),
(1, 7),
(2, 7),
(3, 7);

--
-- Índices para tablas volcadas
--

--
-- Indices de la tabla `chat`
--
ALTER TABLE `chat`
  ADD PRIMARY KEY (`id`);

--
-- Indices de la tabla `clavesmensajes`
--
ALTER TABLE `clavesmensajes`
  ADD PRIMARY KEY (`id`);

--
-- Indices de la tabla `clavesusuario`
--
ALTER TABLE `clavesusuario`
  ADD PRIMARY KEY (`idusuario`,`idclavesmensajes`),
  ADD KEY `clavesusuario_rest2` (`idclavesmensajes`);

--
-- Indices de la tabla `mensaje`
--
ALTER TABLE `mensaje`
  ADD PRIMARY KEY (`id`),
  ADD KEY `chat` (`chat`),
  ADD KEY `clave` (`clave`),
  ADD KEY `emisor` (`emisor`);

--
-- Indices de la tabla `receptoresmensaje`
--
ALTER TABLE `receptoresmensaje`
  ADD PRIMARY KEY (`idmensaje`,`idreceptor`),
  ADD KEY `receptoresmensaje_rest2` (`idreceptor`);

--
-- Indices de la tabla `usuario`
--
ALTER TABLE `usuario`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `nombre` (`nombre`);

--
-- Indices de la tabla `usuarioschat`
--
ALTER TABLE `usuarioschat`
  ADD PRIMARY KEY (`idusuario`,`idchat`),
  ADD KEY `idchat` (`idchat`);

--
-- AUTO_INCREMENT de las tablas volcadas
--

--
-- AUTO_INCREMENT de la tabla `chat`
--
ALTER TABLE `chat`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=8;
--
-- AUTO_INCREMENT de la tabla `clavesmensajes`
--
ALTER TABLE `clavesmensajes`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=4;
--
-- AUTO_INCREMENT de la tabla `mensaje`
--
ALTER TABLE `mensaje`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=13;
--
-- AUTO_INCREMENT de la tabla `usuario`
--
ALTER TABLE `usuario`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=26;
--
-- Restricciones para tablas volcadas
--

--
-- Filtros para la tabla `clavesusuario`
--
ALTER TABLE `clavesusuario`
  ADD CONSTRAINT `clavesusuario_rest1` FOREIGN KEY (`idusuario`) REFERENCES `usuario` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `clavesusuario_rest2` FOREIGN KEY (`idclavesmensajes`) REFERENCES `clavesmensajes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `mensaje`
--
ALTER TABLE `mensaje`
  ADD CONSTRAINT `mensaje_ibfk_1` FOREIGN KEY (`chat`) REFERENCES `chat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `mensaje_ibfk_2` FOREIGN KEY (`clave`) REFERENCES `clavesmensajes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `mensaje_ibfk_3` FOREIGN KEY (`emisor`) REFERENCES `usuario` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `receptoresmensaje`
--
ALTER TABLE `receptoresmensaje`
  ADD CONSTRAINT `receptoresmensaje_rest1` FOREIGN KEY (`idmensaje`) REFERENCES `mensaje` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `receptoresmensaje_rest2` FOREIGN KEY (`idreceptor`) REFERENCES `usuario` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `usuarioschat`
--
ALTER TABLE `usuarioschat`
  ADD CONSTRAINT `usuarioschat_ibfk_1` FOREIGN KEY (`idusuario`) REFERENCES `usuario` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `usuarioschat_ibfk_2` FOREIGN KEY (`idchat`) REFERENCES `chat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
