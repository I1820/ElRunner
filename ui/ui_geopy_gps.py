from geopy.geocoders import Nominatim

# library link: https://github.com/geopy/geopy
# Geocoding
# To geolocate a query to an address and coordinates:
geolocator = Nominatim()
location = geolocator.geocode("175 5th Avenue NYC")
print(location.address)
print(location.latitude, location.longitude)
print(location.raw)

# To find the address corresponding to a set of coordinates:
location = geolocator.reverse("52.509669, 13.376294")
print(location.address)
print(location.latitude, location.longitude)
print(location.raw)
