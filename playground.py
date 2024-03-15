import tkinter as tk
from tkinter import ttk

from pyproj import Transformer


# Define a function to perform the coordinate transformation
def transform_coordinates():
    from_crs = from_crs_combobox.get()
    to_crs = to_crs_combobox.get()
    lat = float(lat_entry.get())
    lon = float(lon_entry.get())

    transformer = Transformer.from_crs(from_crs, to_crs)
    x, y = transformer.transform(lat, lon)

    result_label.config(text=f"Transformed coordinates: {x:.4f}, {y:.4f}")


# Create the main window
root = tk.Tk()
root.title("CRS Transformer")

# Create and place the widgets
tk.Label(root, text="Latitude:").grid(row=0, column=0)
lat_entry = tk.Entry(root)
lat_entry.grid(row=0, column=1)

tk.Label(root, text="Longitude:").grid(row=1, column=0)
lon_entry = tk.Entry(root)
lon_entry.grid(row=1, column=1)

tk.Label(root, text="From CRS:").grid(row=2, column=0)
from_crs_combobox = ttk.Combobox(root, values=["EPSG:4326", "EPSG:3857", "EPSG:32633"])
from_crs_combobox.grid(row=2, column=1)
from_crs_combobox.set("EPSG:4326")

tk.Label(root, text="To CRS:").grid(row=3, column=0)
to_crs_combobox = ttk.Combobox(root, values=["EPSG:4326", "EPSG:3857", "EPSG:32633"])
to_crs_combobox.grid(row=3, column=1)
to_crs_combobox.set("EPSG:32633")

transform_button = tk.Button(root, text="Transform", command=transform_coordinates)
transform_button.grid(row=4, column=0, columnspan=2)

result_label = tk.Label(root, text="")
result_label.grid(row=5, column=0, columnspan=2)

# Start the GUI event loop
root.mainloop()
