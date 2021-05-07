import { IMovie } from '../models/MovieInfo';
import image1 from '../images/image_1.jpeg';
import image2 from '../images/image_2.jpeg';
import image3 from '../images/image_3.jpeg';
import image4 from '../images/image_4.jpeg';
import video1 from '../images/video01.webm';
import video2 from '../images/video02.webm';
import video3 from '../images/video03.webm';

export const cards: IMovie[] = [
  {
    id: '0',
    title: 'Capitan Marvel',
    img:
      'https://upload.wikimedia.org/wikipedia/ru/0/07/Captain_Marvel_film_logo.jpg',
    src: video1,
    info: {
      avalibility: 0.3,
      year: 2000,
      genres: ['Action', 'Sci-Fi'],
      countries: ['USA'],
      rating: 3.4,
      views: 123000,
      length: 123,
      pgRating: 'PG-13',
      description:
        "Carol Danvers becomes one of the universe's most powerful heroes when Earth is caught in the middle of a galactic war between two alien races.",
      photos: [image1, image2, image3, image4],
      videos: [video1, video2, video3],
      imdbRating: 4,
      maxComments: 2,
      comments: [
        {
          id: 1,
          movieid: '0',
          time: 0,
          text: 'hello',
          username: 'Alex',
          avatar: 'https://randomuser.me/api/portraits/men/2.jpg',
        },
        {
          id: 2,
          movieid: '0',
          time: 0,
          text: 'nice movie',
          username: 'Lee',
          avatar: 'https://randomuser.me/api/portraits/women/2.jpg',
        },
      ],
    },
  },
];
