-- phpMyAdmin SQL Dump
-- version 4.4.12
-- http://www.phpmyadmin.net
--
-- Servidor: 127.0.0.1
-- Tiempo de generación: 29-03-2016 a las 17:42:42
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
  `id` int(11) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `chat`
--

INSERT INTO `chat` (`id`) VALUES
(1),
(2),
(3),
(4),
(5);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `claveschat`
--

CREATE TABLE IF NOT EXISTS `claveschat` (
  `id` int(11) NOT NULL,
  `clave` varchar(50) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `claveschat`
--

INSERT INTO `claveschat` (`id`, `clave`) VALUES
(1, 'clave1cifradachat');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `clavesusuario`
--

CREATE TABLE IF NOT EXISTS `clavesusuario` (
  `id` int(11) NOT NULL,
  `clave` varchar(50) NOT NULL,
  `usuario` int(11) NOT NULL,
  `clavechat` int(11) DEFAULT NULL
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `clavesusuario`
--

INSERT INTO `clavesusuario` (`id`, `clave`, `usuario`, `clavechat`) VALUES
(1, 'clave1cifrada', 1, NULL);

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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `mensaje`
--

INSERT INTO `mensaje` (`id`, `texto`, `emisor`, `chat`, `clave`) VALUES
(2, 'Hola que tal?? :)', 1, 5, 1);

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
(2, 15, 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuario`
--

CREATE TABLE IF NOT EXISTS `usuario` (
  `id` int(11) NOT NULL,
  `nombre` varchar(15) NOT NULL,
  `clavepubrsa` varchar(50) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=latin1;

--
-- Volcado de datos para la tabla `usuario`
--

INSERT INTO `usuario` (`id`, `nombre`, `clavepubrsa`) VALUES
(1, 'pepe', 'clave1'),
(13, 'lucia', 'clave1'),
(15, 'maria', 'clave1');

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
(15, 5);

--
-- Índices para tablas volcadas
--

--
-- Indices de la tabla `chat`
--
ALTER TABLE `chat`
  ADD PRIMARY KEY (`id`);

--
-- Indices de la tabla `claveschat`
--
ALTER TABLE `claveschat`
  ADD PRIMARY KEY (`id`);

--
-- Indices de la tabla `clavesusuario`
--
ALTER TABLE `clavesusuario`
  ADD PRIMARY KEY (`id`),
  ADD KEY `usuario` (`usuario`),
  ADD KEY `clavechat` (`clavechat`);

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
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=6;
--
-- AUTO_INCREMENT de la tabla `claveschat`
--
ALTER TABLE `claveschat`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT de la tabla `clavesusuario`
--
ALTER TABLE `clavesusuario`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT de la tabla `mensaje`
--
ALTER TABLE `mensaje`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=3;
--
-- AUTO_INCREMENT de la tabla `usuario`
--
ALTER TABLE `usuario`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=16;
--
-- Restricciones para tablas volcadas
--

--
-- Filtros para la tabla `clavesusuario`
--
ALTER TABLE `clavesusuario`
  ADD CONSTRAINT `clavesusuario_ibfk_1` FOREIGN KEY (`usuario`) REFERENCES `usuario` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `clavesusuario_ibfk_2` FOREIGN KEY (`clavechat`) REFERENCES `claveschat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `mensaje`
--
ALTER TABLE `mensaje`
  ADD CONSTRAINT `mensaje_ibfk_1` FOREIGN KEY (`chat`) REFERENCES `chat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `mensaje_ibfk_2` FOREIGN KEY (`clave`) REFERENCES `claveschat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
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
