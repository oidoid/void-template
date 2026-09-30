<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.10" tiledversion="1.12.2" name="entities" tilewidth="8" tileheight="13" tilecount="2" columns="0" objectalignment="topleft">
 <grid orientation="orthogonal" width="1" height="1"/>
 <tile id="0" type="P1">
  <image source="../atlas/p1.aseprite" width="8" height="13"/>
  <properties>
   <property name="Cel" type="int" propertytype="Cel" value="0"/>
   <property name="Tag" value="P1WalkRight"/>
   <property name="Layer" type="int" propertytype="Layer" value="1"/>
   <property name="Pal" type="int" propertytype="Pal" value="0"/>
   <property name="Stretch" type="bool" value="true"/>
   <property name="Z" type="int" propertytype="Z" value="0"/>
   <property name="ZTop" type="bool" value="false"/>
  </properties>
 </tile>
 <tile id="1" type="Cursor">
  <properties>
   <property name="Cel" type="int" propertytype="Cel" value="0"/>
   <property name="KbdVel" type="int" value="100"/>
   <property name="Layer" type="int" propertytype="Layer" value="3"/>
   <property name="Pal" type="int" propertytype="Pal" value="0"/>
   <property name="Stretch" type="bool" value="true"/>
   <property name="Tag" value="CursorPoint"/>
   <property name="Z" type="int" propertytype="Z" value="0"/>
  </properties>
  <image source="../atlas/cursor.aseprite" width="8" height="13"/>
 </tile>
</tileset>
