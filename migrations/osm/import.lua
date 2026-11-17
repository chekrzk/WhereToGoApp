local schema = 'osm'
local srid = 4326

local pois = osm2pgsql.define_table({
    name = 'pois',
    schema = schema,

    ids = {
        type = 'any',
        id_column = 'osm_id',
        type_column = 'osm_type',
        create_index = 'primary_key'
    },

    columns = {
        { column = 'name',        type = 'text' },
        { column = 'category',    type = 'text',  not_null = true },
        { column = 'subcategory', type = 'text',  not_null = true },
        { column = 'tags',        type = 'jsonb', not_null = true },
        {
            column = 'geom',
            type = 'point',
            projection = srid,
            not_null = true
        }
    },

    indexes = {
        { column = 'geom', method = 'gist' },
        {
            column = { 'category', 'subcategory' },
            method = 'btree'
        }
    }
})

local areas = osm2pgsql.define_table({
    name = 'areas',
    schema = schema,

    ids = {
        type = 'any',
        id_column = 'osm_id',
        type_column = 'osm_type',
        create_index = 'primary_key'
    },

    columns = {
        { column = 'name',        type = 'text' },
        { column = 'category',    type = 'text',  not_null = true },
        { column = 'subcategory', type = 'text',  not_null = true },
        { column = 'tags',        type = 'jsonb', not_null = true },
        {
            column = 'geom',
            type = 'geometry',
            projection = srid,
            not_null = true
        }
    },

    indexes = {
        { column = 'geom', method = 'gist' },
        {
            column = { 'category', 'subcategory' },
            method = 'btree'
        }
    }
})

local included_landuse = {
    residential = true,
    commercial = true,
    retail = true,
    industrial = true,
    recreation_ground = true,
    forest = true,
    meadow = true,
    grass = true
}

local included_leisure = {
    park = true,
    garden = true,
    nature_reserve = true
}

local included_natural = {
    wood = true,
    grassland = true
}

local function classify_poi(tags)
    if tags.shop then
        return 'shop', tags.shop
    end

    if tags.amenity then
        return 'amenity', tags.amenity
    end

    if tags.tourism then
        return 'tourism', tags.tourism
    end

    return nil, nil
end

local function classify_area(tags)
    if tags.leisure and included_leisure[tags.leisure] then
        return 'leisure', tags.leisure
    end

    if tags.landuse and included_landuse[tags.landuse] then
        return 'landuse', tags.landuse
    end

    if tags.natural and included_natural[tags.natural] then
        return 'natural', tags.natural
    end

    return nil, nil
end

local function insert_poi(object, geom)
    local category, subcategory = classify_poi(object.tags)

    if not category then
        return
    end

    pois:insert({
        name = object.tags.name,
        category = category,
        subcategory = subcategory,
        tags = object.tags,
        geom = geom
    })
end

local function insert_area(object, geom)
    local category, subcategory = classify_area(object.tags)

    if not category then
        return
    end

    areas:insert({
        name = object.tags.name,
        category = category,
        subcategory = subcategory,
        tags = object.tags,
        geom = geom
    })
end

function osm2pgsql.process_node(object)
    insert_poi(object, object:as_point())
end

function osm2pgsql.process_way(object)
    if not object.is_closed or object.tags.area == 'no' then
        return
    end

    local polygon = object:as_polygon()

    insert_poi(object, polygon:centroid())
    insert_area(object, polygon)
end

function osm2pgsql.process_relation(object)
    if object.tags.type ~= 'multipolygon' then
        return
    end

    local multipolygon = object:as_multipolygon()

    insert_poi(object, multipolygon:centroid())
    insert_area(object, multipolygon)
end