# awegeo
`awegeo` is a simple command line tool that converts a KML/KMZ file of
the [AWARE](http://www.wulca-waterlca.org/aware.html) method into a
[GeoJSON](https://tools.ietf.org/html/rfc7946) file.


```xml
<Placemark id="ID_00000">
  <name>0</name>
  <Snippet></Snippet>
  <description><![CDATA[
<html ...>
<head>
    ...
</head>
<body style="...">
 <table style="...">
  <tr style="...">
    <td>0</td>
  </tr>
  <tr>
    <td>
     <table style="...">
         <tr>
             <td>FID</td>
             <td>0</td>
         </tr>
         <tr bgcolor="#D4E4F3">
             <td>Consumption_m3</td>
             <td>0.0e+000</td>
         </tr>
         <tr>
             <td>Area_m2</td>
             <td>3.4e+008</td>
         </tr>
         ...
     </table>
    </td>
  </tr>
 </table>
</body>
</html>]]></description>
  <styleUrl>#PolyStyle00</styleUrl>
  <MultiGeometry>
    <Polygon>
      <extrude>0</extrude>
      <altitudeMode>clampToGround</altitudeMode>
      <outerBoundaryIs>
        <LinearRing>
          <coordinates>
            -38.00002600000005,83.49997399999998,0
            -38.00002600000005,83.99997399999998,0
            -38.50002600000005,83.99997399999998,0
            -38.50002600000005,83.49997399999998,0
            -38.00002600000005,83.49997399999998,0
          </coordinates>
        </LinearRing>
      </outerBoundaryIs>
    </Polygon>
  </MultiGeometry>
</Placemark>
```

`<tr><td>([^<]*)<\/td><td>([^<]*)<\/td><\/tr>`