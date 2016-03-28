-- phpMyAdmin SQL Dump
-- version 4.2.12deb2
-- http://www.phpmyadmin.net
--
-- Servidor: localhost
-- Tiempo de generación: 28-03-2016 a las 17:50:50
-- Versión del servidor: 5.5.44-0+deb8u1
-- Versión de PHP: 5.6.9-0+deb8u1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- Base de datos: `sds`
--

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `chat`
--

CREATE TABLE IF NOT EXISTS `chat` (
`id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `claveschat`
--

CREATE TABLE IF NOT EXISTS `claveschat` (
`id` int(11) NOT NULL,
  `clave` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `clavesusuario`
--

CREATE TABLE IF NOT EXISTS `clavesusuario` (
`id` int(11) NOT NULL,
  `clave` varchar(50) NOT NULL,
  `usuario` int(11) NOT NULL,
  `clavechat` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `mensaje`
--

CREATE TABLE IF NOT EXISTS `mensaje` (
`id` int(11) NOT NULL,
  `texto` varchar(1000) NOT NULL,
  `chat` int(11) NOT NULL,
  `clave` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuario`
--

CREATE TABLE IF NOT EXISTS `usuario` (
`id` int(11) NOT NULL,
  `nombre` varchar(15) NOT NULL,
  `clavepubrsa` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuariochat`
--

CREATE TABLE IF NOT EXISTS `usuariochat` (
  `idusuario` int(11) NOT NULL,
  `idchat` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

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
 ADD PRIMARY KEY (`id`), ADD KEY `usuario` (`usuario`), ADD KEY `clavechat` (`clavechat`);

--
-- Indices de la tabla `mensaje`
--
ALTER TABLE `mensaje`
 ADD PRIMARY KEY (`id`), ADD KEY `chat` (`chat`), ADD KEY `clave` (`clave`);

--
-- Indices de la tabla `usuario`
--
ALTER TABLE `usuario`
 ADD PRIMARY KEY (`id`), ADD UNIQUE KEY `nombre` (`nombre`);

--
-- Indices de la tabla `usuariochat`
--
ALTER TABLE `usuariochat`
 ADD PRIMARY KEY (`idusuario`,`idchat`), ADD KEY `idchat` (`idchat`);

--
-- AUTO_INCREMENT de las tablas volcadas
--

--
-- AUTO_INCREMENT de la tabla `chat`
--
ALTER TABLE `chat`
MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT de la tabla `claveschat`
--
ALTER TABLE `claveschat`
MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT de la tabla `clavesusuario`
--
ALTER TABLE `clavesusuario`
MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT de la tabla `mensaje`
--
ALTER TABLE `mensaje`
MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT de la tabla `usuario`
--
ALTER TABLE `usuario`
MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- Restricciones para tablas volcadas
--

--
-- Filtros para la tabla `clavesusuario`
--
ALTER TABLE `clavesusuario`
ADD CONSTRAINT `clavesusuario_ibfk_2` FOREIGN KEY (`clavechat`) REFERENCES `claveschat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT `clavesusuario_ibfk_1` FOREIGN KEY (`usuario`) REFERENCES `usuario` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `mensaje`
--
ALTER TABLE `mensaje`
ADD CONSTRAINT `mensaje_ibfk_2` FOREIGN KEY (`clave`) REFERENCES `claveschat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT `mensaje_ibfk_1` FOREIGN KEY (`chat`) REFERENCES `chat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `usuariochat`
--
ALTER TABLE `usuariochat`
ADD CONSTRAINT `usuariochat_ibfk_2` FOREIGN KEY (`idchat`) REFERENCES `chat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT `usuariochat_ibfk_1` FOREIGN KEY (`idusuario`) REFERENCES `usuario` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
